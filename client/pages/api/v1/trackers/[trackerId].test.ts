import { mockReqRes } from 'tests/api/common';
import getTrackerByIdHandler from './[trackerId]';

describe('getTrackerByIdHandler', () => {
  test('it returns a 405 if called with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getTrackerByIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test.todo('it returns a 401');
  test.todo('it returns a 404');
  test.todo('it successfully returns a tracker');
});
