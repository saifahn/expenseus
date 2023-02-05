import { mockReqRes } from 'tests/api/common';
import getTxnsByUserIdBetweenDatesHandler from './range';

describe('getTxnsByUserIdBetweenDates handler', () => {
  test('it returns a 405 for non-GET requests', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByUserIdBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 400 if there is no from or to query or they cannot be parsed to numbers', async () => {
    var { req, res } = mockReqRes('GET');
    await getTxnsByUserIdBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: 'not-a-number',
      to: '12345678',
    };
    await getTxnsByUserIdBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
    };
    await getTxnsByUserIdBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);
  });
});
