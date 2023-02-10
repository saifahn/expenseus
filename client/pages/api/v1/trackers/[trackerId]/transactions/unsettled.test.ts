import { mockReqRes } from 'tests/api/common';
import getUnsettledTxnsByTrackerHandler from './unsettled';

describe('getUnsettledTxnsByTrackerHandler', () => {
  test('it returns 405 on a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getUnsettledTxnsByTrackerHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
});
