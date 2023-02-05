import { mockReqRes } from 'tests/api/common';
import getTxnsByUserIdBetweenDatesHandler from './range';

describe('getTxnsByUserIdBetweenDates handler', () => {
  test('it returns a 405 for non-GET requests', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByUserIdBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
});
