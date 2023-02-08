import { mockReqRes } from 'tests/api/common';
import bySharedTxnIdHandler from './[transactionId]';

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
});
