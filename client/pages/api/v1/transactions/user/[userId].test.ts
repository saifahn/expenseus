import {
  tablePartitionKey,
  tableSortKey,
  gsi1PartitionKey,
  gsi1SortKey,
} from 'ddb/schema';
import { makeTxnRepository } from 'ddb/txns';
import * as nextAuth from 'next-auth';
import { assertEqualTxnDetails, mockReqRes } from 'tests/api/common';
import { txnRepoFnsMock } from '../../transactions.test';
import txnByUserIdHandler from './[userId]';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);

jest.mock('next-auth');
const nextAuthMocked = jest.mocked(nextAuth);

describe('txnByUserId handler', () => {
  test('it returns a 405 for a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await txnByUserIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a list of transactions for the user', async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMocked.getServerSession.mockImplementationOnce(async () => {
      return {
        user: {
          email: 'test-user',
        },
      };
    });
    req.query.userId = 'test-user';
    const mockTxnItem = {
      [tablePartitionKey]: 'user#test-user',
      [tableSortKey]: 'txn#test-txn',
      [gsi1PartitionKey]: 'user#test-user',
      [gsi1SortKey]: 'txn#12345678#test-txn',
      EntityType: 'transaction' as const,
      ID: 'test-txn',
      UserID: 'test-user',
      Date: 12345678,
      Amount: 9275,
      Location: 'somewhere',
      Category: 'unspecified.unspecified' as const,
      Details: '',
    };
    txnsRepo.mockImplementationOnce(() => ({
      ...txnRepoFnsMock,
      getTxnsByUserId: jest.fn(async () => [mockTxnItem]),
    }));
    await txnByUserIdHandler(req, res);

    expect(res.statusCode).toBe(200);
    const result = res._getJSONData();
    expect(result).toHaveLength(1);
    assertEqualTxnDetails(result[0], mockTxnItem);
  });

  test("it returns a 403 when a user attempts to retrieve another user's txns", async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMocked.getServerSession.mockImplementationOnce(async () => {
      return {
        user: {
          email: 'a-different-user',
        },
      };
    });
    req.query.userId = 'test-user';
    await txnByUserIdHandler(req, res);

    expect(res.statusCode).toBe(403);
  });
});
