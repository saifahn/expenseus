import * as nextAuth from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import createTrackerHandler from './trackers';

const nextAuthMock = jest.mocked(nextAuth);

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

  test('returns a 401 with no valid session', async () => {
    const { req, res } = mockReqRes('POST');
    req._setBody({ users: ['test-user', 'test-user-2'], name: 'test-tracker' });
    nextAuthMock.getServerSession.mockResolvedValueOnce(null);
    await createTrackerHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  it("returns a 403 error when a user attempts to create a tracker that doesn't include them", async () => {
    const { req, res } = mockReqRes('POST');
    req._setBody({ users: ['test-user', 'test-user-2'], name: 'test-tracker' });
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: {
        email: 'different-user',
      },
    });
    await createTrackerHandler(req, res);

    expect(res.statusCode).toBe(403);
  });
});
