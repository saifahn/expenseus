import { mockReqRes } from 'tests/api/common';
import bySharedTxnIdHandler from './[transactionId]';

describe('bySharedTxnIdHandler', () => {
  test('it returns a 405 if the method is not PUT or DELETE', async () => {
    const { req, res } = mockReqRes('GET');
    await bySharedTxnIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
});
