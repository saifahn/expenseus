import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import getTxnsByTrackerBetweenDatesHandler from './range';

const sessionMock = jest.mocked(getServerSession);

describe('getTxnsByTrackerBetweenDatesHandler', () => {
  test('it returns a 405 for non-GET requests', async () => {
    const { req, res } = mockReqRes('POST');
    await getTxnsByTrackerBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 400 if there is no from or to query or they cannot be parsed to numbers', async () => {
    var { req, res } = mockReqRes('GET');
    await getTxnsByTrackerBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: 'not-a-number',
      to: '12345678',
    };
    await getTxnsByTrackerBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);

    var { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
    };
    await getTxnsByTrackerBetweenDatesHandler(req, res);
    expect(res.statusCode).toBe(400);
  });

  test('it returns a 401 if there is no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    req.query = {
      from: '12345678',
      to: '23456789',
    };
    sessionMock.mockResolvedValueOnce(null);
    await getTxnsByTrackerBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(401);
  });
});
