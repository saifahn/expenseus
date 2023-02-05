import { mockReqRes } from 'tests/api/common';
import settleTxnsHandler from './settle';

describe('settleTxnsHandler', () => {
  test('returns a 405 if called with a non-POST method', async () => {
    const { req, res } = mockReqRes('GET');
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('returns a 400 if the input is invalid', async () => {
    const { req, res } = mockReqRes('POST');
    req._setBody({ something: 'invalid' });
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(400);
  });
});
