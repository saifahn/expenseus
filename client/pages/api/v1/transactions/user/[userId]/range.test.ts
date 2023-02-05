import { makeTxnRepository } from 'ddb/txns';
import * as nextAuth from 'next-auth';
import { txnRepoFnsMock } from 'pages/api/v1/transactions.test';
import { assertEqualTxnDetails, mockReqRes } from 'tests/api/common';
import { mockTxnItem } from 'tests/api/doubles';
import getTxnsByUserIdBetweenDatesHandler from './range';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);
const nextAuthMock = jest.mocked(nextAuth);

describe('getTxnsByUserIdBetweenDates handler', () => {
  test('it returns a 405 for non-GET requests', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByUserIdBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 400 if there is no from or to query or they cannot be parsed to numbers', async () => {
    var { req, res } = mockReqRes('GET');
    await getTxnsByUserIdBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: 'not-a-number',
      to: '12345678',
    };
    await getTxnsByUserIdBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
    };
    await getTxnsByUserIdBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);
  });

  test('it returns a 401 if there is no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
      to: '23456789',
    };
    nextAuthMock.getServerSession.mockResolvedValueOnce(null);
    await getTxnsByUserIdBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('it returns a list of txns based on the result from the store', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
      to: '23456789',
    };
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: {
        email: 'test-user',
      },
    });
    const getBetweenDatesMock = jest.fn(async () => [mockTxnItem]);
    txnsRepo.mockImplementationOnce(() => ({
      ...txnRepoFnsMock,
      getBetweenDates: getBetweenDatesMock,
    }));
    await getTxnsByUserIdBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(200);
    expect(getBetweenDatesMock).toHaveBeenCalledWith({
      from: 12345678,
      to: 23456789,
      userId: 'test-user',
    });
    const result = res._getJSONData();
    expect(result).toHaveLength(1);
    assertEqualTxnDetails(result[0], mockTxnItem);
  });
});
