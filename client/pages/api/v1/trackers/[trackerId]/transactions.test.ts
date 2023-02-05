import { mockReqRes } from 'tests/api/common';
import txnsByTrackerHandler from './transactions';

describe('txnsByTrackerHandler', () => {
  test('returns a 405 if called with a non-POST or GET method', async () => {
    const { req, res } = mockReqRes('DELETE');
    await txnsByTrackerHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
  describe('GET all', () => {
    test.todo('');
  });
});
