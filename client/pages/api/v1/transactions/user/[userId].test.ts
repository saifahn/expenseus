import { makeTxnRepository } from 'ddb/txns';
import * as nextAuth from 'next-auth';
import { assertEqualTxnDetails, mockReqRes } from 'tests/api/common';
import { mockTxnItem } from 'tests/api/doubles';
import { txnRepoFnsMock } from '../../transactions.test';
import getTxnsByUserIdHandler from './[userId]';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);
const nextAuthMock = jest.mocked(nextAuth);

describe('txnByUserId handler', () => {
  test('it returns a 405 for a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByUserIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a list of transactions for the user', async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: {
        email: 'test-user',
      },
    });
    req.query.userId = 'test-user';
    txnsRepo.mockImplementationOnce(() => ({
      ...txnRepoFnsMock,
      getTxnsByUserId: jest.fn(async () => [mockTxnItem]),
    }));
    await getTxnsByUserIdHandler(req, res);

    expect(res.statusCode).toBe(200);
    const result = res._getJSONData();
    expect(result).toHaveLength(1);
    assertEqualTxnDetails(result[0], mockTxnItem);
  });

  test("it returns a 403 when a user attempts to retrieve another user's txns", async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: {
        email: 'a-different-user',
      },
    });
    req.query.userId = 'test-user';
    await getTxnsByUserIdHandler(req, res);

    expect(res.statusCode).toBe(403);
  });

  test('it returns a 401 when there is no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMock.getServerSession.mockResolvedValueOnce(null);
    await getTxnsByUserIdHandler(req, res);

    expect(res.statusCode).toBe(401);
  });
});
