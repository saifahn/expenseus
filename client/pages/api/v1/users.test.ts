import usersHandler from './users';
import { makeUserRepository } from 'ddb/users';
import { mockReqRes } from 'tests/api/common';
import { testUserItem, userRepoFnsMock } from 'tests/api/doubles';
import { getServerSession } from 'next-auth';

jest.mock('ddb/users');
const usersRepo = jest.mocked(makeUserRepository);
const serverSessionMock = jest.mocked(getServerSession);

describe('/api/v1/users API endpoint', () => {
  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('should return a list of users when they are returned from the store', async () => {
    usersRepo.mockReturnValueOnce({
      ...userRepoFnsMock,
      getAllUsers: jest.fn(async () => [testUserItem]),
    });

    serverSessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    const { req, res } = mockReqRes('GET');
    await usersHandler(req, res);

    expect(res.statusCode).toBe(200);
    expect(res.getHeaders()).toEqual({ 'content-type': 'application/json' });

    const expected = {
      username: 'testUser',
      name: 'Test User',
      id: 'test-user',
    };
    const result = res._getJSONData();
    expect(result).toHaveLength(1);
    expect(result).toContainEqual(expected);
  });

  it('returns a 405 error when called with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    serverSessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    await usersHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  it('returns a 401 error when there is no valid session', async () => {
    // mock the session to return invalid
    const { req, res } = mockReqRes('GET');
    serverSessionMock.mockResolvedValueOnce(null);
    await usersHandler(req, res);

    expect(res.statusCode).toBe(401);
  });
});
