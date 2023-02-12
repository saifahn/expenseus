import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import getAllTxnsByUserBetweenDatesHandler from './all';

const sessionMock = jest.mocked(getServerSession);

describe('getAllTxnsByUserBetweenDatesHandler', () => {
  test('it returns a 405 with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns 401 with no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test("it returns 403 when trying to get someone else's txns", async () => {
    const { req, res } = mockReqRes('GET');
    req.query.userId = 'different-user';
    sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    await getAllTxnsByUserBetweenDatesHandler(req, res);

    expect(res.statusCode).toBe(403);
  });
  test.todo('it returns 400 with incorrect from, to');
  test.todo('it returns a list of txns and shared txns for the user');
});
