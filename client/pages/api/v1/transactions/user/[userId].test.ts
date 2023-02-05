import { mockReqRes } from 'tests/api/common';
import txnByUserIdHandler from './[userId]';

describe('txnByUserId handler', () => {
  test('it returns a 405 for a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await txnByUserIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
});
