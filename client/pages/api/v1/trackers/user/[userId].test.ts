import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import getTrackersByUserHandler from './[userId]';

const sessionMock = jest.mocked(getServerSession);

describe('getTrackersByUserHandler', () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test('it returns a 405 when called with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getTrackersByUserHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 401 when called with no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    await getTrackersByUserHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('it returns a 403 when called for a user different to logged in user', async () => {
    const { req, res } = mockReqRes('GET');
    sessionMock.mockResolvedValueOnce({ user: { email: 'different-user' } });
    req.query = { userId: 'test-user' };
    await getTrackersByUserHandler(req, res);

    expect(res.statusCode).toBe(403);
  });
  test.todo('it returns a list of trackers for a user successfully');
});
