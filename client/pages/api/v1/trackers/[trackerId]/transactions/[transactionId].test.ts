import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import bySharedTxnIdHandler from './[transactionId]';

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
  });
});
