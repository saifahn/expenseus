import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import { getServerSession } from 'next-auth';
import { assertEqualSharedTxnDetails, mockReqRes } from 'tests/api/common';
import { mockSharedTxnItem, sharedTxnRepoFnsMock } from 'tests/api/doubles';
import getTxnsByTrackerBetweenDatesHandler from './range';

jest.mock('ddb/sharedTxns');
const sharedTxnsRepo = jest.mocked(makeSharedTxnRepository);
const sessionMock = jest.mocked(getServerSession);

describe('getTxnsByTrackerBetweenDatesHandler', () => {
  test('it returns a 405 for non-GET requests', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByTrackerBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 400 if there is no from or to query or they cannot be parsed to numbers', async () => {
    var { req, res } = mockReqRes('GET');
    await getTxnsByTrackerBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: 'not-a-number',
      to: '12345678',
    };
    await getTxnsByTrackerBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
    };
    await getTxnsByTrackerBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);
  });

  test('it returns a 401 if there is no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
      to: '23456789',
      trackerId: 'test-tracker',
    };
    sessionMock.mockResolvedValueOnce(null);
    await getTxnsByTrackerBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('it returns a list of shared txns based on the results from the store', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
      to: '23456789',
      trackerId: 'test-tracker',
    };
    sessionMock.mockResolvedValueOnce({
      user: {
        email: 'test-user',
      },
    });
    const getBetweenDatesMock = jest.fn(async () => [mockSharedTxnItem]);
    sharedTxnsRepo.mockReturnValueOnce({
      ...sharedTxnRepoFnsMock,
      getTxnsByTrackerBetweenDates: getBetweenDatesMock,
    });
    await getTxnsByTrackerBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(200);
    expect(getBetweenDatesMock).toHaveBeenCalledWith({
      from: 12345678,
      to: 23456789,
      tracker: 'test-tracker',
    });
    const result = res._getJSONData();
    expect(result).toHaveLength(1);
    assertEqualSharedTxnDetails(result[0], mockSharedTxnItem);
  });
});
