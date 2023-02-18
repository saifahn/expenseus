import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import { makeTxnRepository } from 'ddb/txns';
import { getServerSession } from 'next-auth';
import { txnRepoFnsMock } from 'pages/api/v1/transactions.test';
import {
  assertEqualSharedTxnDetails,
  assertEqualTxnDetails,
  mockReqRes,
} from 'tests/api/common';
import {
  mockSharedTxnItem,
  mockTxnItem,
  sharedTxnRepoFnsMock,
} from 'tests/api/doubles';
import getAllTxnsByUserBetweenDatesHandler from './all';

jest.mock('ddb/txns');
jest.mock('ddb/sharedTxns');
const txnRepo = jest.mocked(makeTxnRepository);
const sharedTxnRepo = jest.mocked(makeSharedTxnRepository);
const sessionMock = jest.mocked(getServerSession);

describe('getAllTxnsByUserBetweenDatesHandler', () => {
  test('it returns a 405 with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns 400 with incorrect from, to', async () => {
    const { req, res } = mockReqRes('GET');
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(400);
  });

  test('it returns 401 with no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '1234',
      to: '2345',
    };
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test("it returns 403 when trying to get someone else's txns", async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '1234',
      to: '2345',
      userId: 'different-user',
    };
    sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(403);
  });

  test('it returns a list of txns and shared txns for the user', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '1234',
      to: '2345',
    };
    req.query.userId = 'test-user';
    sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    // mock the sharedTxnRepo
    const getAllSharedTxnsMock = jest.fn(async () => [mockSharedTxnItem]);
    sharedTxnRepo.mockReturnValueOnce({
      ...sharedTxnRepoFnsMock,
      getSharedTxnsByUserBetweenDates: getAllSharedTxnsMock,
    });
    // mock the txnRepo
    const getTxnsBetweenDatesMock = jest.fn(async () => [mockTxnItem]);
    txnRepo.mockReturnValueOnce({
      ...txnRepoFnsMock,
      getBetweenDates: getTxnsBetweenDatesMock,
    });
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(200);
    const data = res._getJSONData();
    assertEqualSharedTxnDetails(data.sharedTransactions[0], mockSharedTxnItem);
    assertEqualTxnDetails(data.transactions[0], mockTxnItem);
  });
});
