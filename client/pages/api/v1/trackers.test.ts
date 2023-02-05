import { mockReqRes } from 'tests/api/common';
import createTrackerHandler from './trackers';

describe('createTrackerHandler', () => {
  test('returns a 405 with a non-POST request', async () => {
    const { req, res } = mockReqRes('GET');
    await createTrackerHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('returns a 400 with an invalid input', async () => {
    const { req, res } = mockReqRes('POST');
    req._setBody({ invalid: 'input' });
    await createTrackerHandler(req, res);

    expect(res.statusCode).toBe(400);
  });
});
