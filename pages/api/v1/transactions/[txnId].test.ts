import { makeTxnRepository } from 'ddb/txns';
import { assertEqualTxnDetails, mockReqRes } from 'tests/api/common';
import { Transaction } from 'types/Transaction';
import { txnRepoFnsMock } from '../transactions.test';
import byTxnIdHandler from './[txnId]';
import * as nextAuth from 'next-auth';
import { mockTxnItem } from 'tests/api/doubles';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);
const nextAuthMock = jest.mocked(nextAuth);

describe('byTxnIdHandler', () => {
  test('a request with no valid session returns a 401', async () => {
    const { req, res } = mockReqRes('GET');
    nextAuthMock.getServerSession.mockResolvedValueOnce(null);
    await byTxnIdHandler(req, res);
    expect(res.statusCode).toBe(401);
  });

  describe('GET txns by ID', () => {
    test('a txn is successfully retrieved for a valid ID', async () => {
      const { req, res } = mockReqRes('GET');
      req.query.txnId = 'test-txn';
      nextAuthMock.getServerSession.mockResolvedValueOnce({
        user: {
          email: 'test-user',
        },
      });
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
      assertEqualTxnDetails(res._getJSONData(), mockTxnItem);
    });
  });

  describe('PUT - update txn', () => {
    test('a txn can be updated correctly', async () => {
      const { req, res } = mockReqRes('PUT');
      req.query.txnId = 'test-txn';
      nextAuthMock.getServerSession.mockResolvedValueOnce({
        user: {
          email: 'test-user',
        },
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
      req.body = JSON.stringify(updatedTxn);
      txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
      await byTxnIdHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(txnRepoFnsMock.updateTxn).toHaveBeenCalledWith(updatedTxn);
    });
  });

  test('a 400 is returned if the input is incorrect', async () => {
    const { req, res } = mockReqRes('PUT');
    req.query.txnId = 'test-txn';
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: {
        email: 'test-user',
      },
    });
    const updatedTxn = {
      id: 'test-txn',
      something: 'is',
      totally: 'up',
      with: 'this',
    };
    req.body = JSON.stringify(updatedTxn);
    txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
    await byTxnIdHandler(req, res);

    expect(res.statusCode).toBe(400);
  });

  test("a 403 is returned if a user tries to update a txn they're not part of", async () => {
    const { req, res } = mockReqRes('PUT');
    req.query.txnId = 'test-txn';
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: {
        email: 'different-user',
      },
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
    req.body = JSON.stringify(updatedTxn);
    txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
    await byTxnIdHandler(req, res);

    expect(res.statusCode).toBe(403);
  });

  describe('DELETE - delete txn', () => {
    test('a txn can be deleted correctly', async () => {
      const { req, res } = mockReqRes('DELETE');
      req.query.txnId = 'test-txn';
      nextAuthMock.getServerSession.mockResolvedValueOnce({
        user: {
          email: 'test-user',
        },
      });
      txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
      await byTxnIdHandler(req, res);

      expect(txnRepoFnsMock.deleteTxn).toHaveBeenCalledWith({
        txnId: 'test-txn',
        userId: 'test-user',
      });
      expect(res.statusCode).toBe(202);
    });
  });
});
