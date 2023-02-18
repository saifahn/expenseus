import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { sharedTxnRepoFnsMock } from 'tests/api/doubles';
import bySharedTxnIdHandler from './[transactionId]';

jest.mock('ddb/sharedTxns');
const sharedTxnRepo = jest.mocked(makeSharedTxnRepository);
const sessionMock = jest.mocked(getServerSession);

describe('bySharedTxnIdHandler', () => {
  test('it returns a 405 if the method is not PUT or DELETE', async () => {
    const { req, res } = mockReqRes('GET');
    await bySharedTxnIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 401 if called with no valid session', async () => {
    const { req, res } = mockReqRes('PUT');
    await bySharedTxnIdHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  describe('PUT - update shared txn', () => {
    test('it returns a 400 with an invalid input', async () => {
      const { req, res } = mockReqRes('PUT');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const updateSharedTxnInput = {
        invalid: 'input',
      };
      req._setBody(updateSharedTxnInput);
      await bySharedTxnIdHandler(req, res);

      expect(res.statusCode).toBe(400);
    });

    test('it successfully updates a shared txn', async () => {
      const { req, res } = mockReqRes('PUT');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const updateSharedTxnInput = {
        date: 123456790,
        amount: 3646,
        location: 'maruetsu',
        category: 'food.groceries',
        participants: ['test-user', 'test-user-2'],
        payer: 'test-user',
        details: 'something up',
      };
      req.query = {
        transactionId: 'test-shared-txn',
        trackerId: 'test-tracker',
      };
      req._setBody(updateSharedTxnInput);
      sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
      await bySharedTxnIdHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(sharedTxnRepoFnsMock.updateSharedTxn).toHaveBeenCalledWith({
        ...updateSharedTxnInput,
        id: 'test-shared-txn',
        tracker: 'test-tracker',
      });
    });

    test.todo(
      'it returns a 403 when trying to update a shared txn without the session user as a participant',
    );
  });

  describe('DELETE - delete shared txn', () => {
    test('it returns a 400 when given an incorrect input', async () => {
      const { req, res } = mockReqRes('DELETE');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      req.query = {
        transactionId: 'test-shared-txn',
        trackerId: 'test-tracker',
      };
      await bySharedTxnIdHandler(req, res);

      expect(res.statusCode).toBe(400);
    });

    test('it successfully deletes a shared txn', async () => {
      const { req, res } = mockReqRes('DELETE');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      req.query = {
        transactionId: 'test-shared-txn',
        trackerId: 'test-tracker',
      };
      const participants = ['test-user', 'test-user-2'];
      req._setBody({ participants });
      sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
      await bySharedTxnIdHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(sharedTxnRepoFnsMock.deleteSharedTxn).toHaveBeenCalledWith({
        participants,
        txnId: 'test-shared-txn',
        tracker: 'test-tracker',
      });
    });

    test.todo(
      "it returns a 404 when trying to delete a txn that doesn't exist",
    );

    test.todo(
      'it returns 403 when trying to delete a txn without the session user as a participant',
    );
  });
});
