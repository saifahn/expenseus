import { mockReqRes } from 'tests/api/common';
import getTxnsByTrackerBetweenDatesHandler from './range';

describe('getTxnsByTrackerBetweenDatesHandler', () => {
  test('it returns a 405 for non-GET requests', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByTrackerBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
});
