import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';
import { makeTxnRepository, TxnItem } from 'ddb/txns';
import { mockReqRes } from 'tests/api/common';
import { Transaction } from 'types/Transaction';
import { txnRepoFnsMock } from '../transactions.test';
import byTxnIdHandler from './[txnId]';
import * as nextAuth from 'next-auth';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);

jest.mock('next-auth');
const nextAuthMocked = jest.mocked(nextAuth);

// helper function to assert details from a txn match txn item
function assertEqualDetails(txn: Transaction, txnItem: TxnItem) {
  expect(txn).toEqual(
    expect.objectContaining({
      userId: txnItem.UserID,
      location: txnItem.Location,
      amount: txnItem.Amount,
      date: txnItem.Date,
      category: txnItem.Category,
      details: txnItem.Details,
    }),
  );
}

describe('byTxnIdHandler', () => {
  test('a request with no valid session returns a 401', async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMocked.getServerSession.mockImplementationOnce(async () => null);
    await byTxnIdHandler(req, res);
    expect(res.statusCode).toBe(401);
  });

  describe('GET txns by ID', () => {
    test('a txn is successfully retrieved for a valid ID', async () => {
      const { req, res } = mockReqRes('GET');
      req.query.txnId = 'test-txn';
      nextAuthMocked.getServerSession.mockImplementationOnce(async () => {
        return {
          user: {
            email: 'test-user',
          },
        };
      });
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
      const getTxnMock = jest.fn(async () => mockTxnItem);
      txnsRepo.mockImplementationOnce(() => ({
        ...txnRepoFnsMock,
        getTxn: getTxnMock,
      }));
      await byTxnIdHandler(req, res);

      expect(res.statusCode).toBe(200);
      expect(getTxnMock).toHaveBeenCalledWith({
        txnId: 'test-txn',
        userId: 'test-user',
      });
      assertEqualDetails(res._getJSONData(), mockTxnItem);
    });
  });

  describe('PUT - update txn', () => {
    test('a txn can be updated correctly', async () => {
      const { req, res } = mockReqRes('PUT');
      req.query.txnId = 'test-txn';
      nextAuthMocked.getServerSession.mockImplementationOnce(async () => {
        return {
          user: {
            email: 'test-user',
          },
        };
      });
      const updatedTxn: Transaction = {
        id: 'test-txn',
        userId: 'test-user',
        date: 12345678,
        amount: 5000,
        location: 'hair cut',
        category: 'beauty.cosmetics',
        details: '',
      };
      req._setBody(updatedTxn);
      txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
      await byTxnIdHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(txnRepoFnsMock.updateTxn).toHaveBeenCalledWith(updatedTxn);
    });
  });

  test('a 400 is returned if the input is incorrect', async () => {
    const { req, res } = mockReqRes('PUT');
    req.query.txnId = 'test-txn';
    nextAuthMocked.getServerSession.mockImplementationOnce(async () => {
      return {
        user: {
          email: 'test-user',
        },
      };
    });
    const updatedTxn = {
      id: 'test-txn',
      something: 'is',
      totally: 'up',
      with: 'this',
    };
    req._setBody(updatedTxn);
    txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
    await byTxnIdHandler(req, res);

    expect(res.statusCode).toBe(400);
  });

  test("a 403 is returned if a user tries to update a txn they're not part of", async () => {
    const { req, res } = mockReqRes('PUT');
    req.query.txnId = 'test-txn';
    nextAuthMocked.getServerSession.mockImplementationOnce(async () => {
      return {
        user: {
          email: 'different-user',
        },
      };
    });
    const updatedTxn: Transaction = {
      id: 'test-txn',
      userId: 'test-user',
      date: 12345678,
      amount: 5000,
      location: 'hair cut',
      category: 'beauty.cosmetics',
      details: '',
    };
    req._setBody(updatedTxn);
    txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
    await byTxnIdHandler(req, res);

    expect(res.statusCode).toBe(403);
  });
});
