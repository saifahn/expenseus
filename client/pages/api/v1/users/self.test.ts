import { makeUserRepository } from 'ddb/users';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { testUserItem, userRepoFnsMock } from 'tests/api/doubles';
import getSelfHandler from './self';

jest.mock('ddb/users');
const userRepo = jest.mocked(makeUserRepository);
const serverSession = jest.mocked(getServerSession);

describe('getSelfHandler', () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test('it returns a 405 for a non-GET call', async () => {
    const { req, res } = mockReqRes('POST');
    await getSelfHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 401 for a non-valid session', async () => {
    const { req, res } = mockReqRes('GET');
    await getSelfHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('it returns user information for the logged in user', async () => {
    const { req, res } = mockReqRes('GET');
    const getUserMock = jest.fn(async () => testUserItem);
    userRepo.mockReturnValueOnce({
      ...userRepoFnsMock,
      getUser: getUserMock,
    });
    serverSession.mockResolvedValueOnce({ user: { email: 'test-user' } });
    await getSelfHandler(req, res);

    expect(getUserMock).toHaveBeenCalledWith('test-user');
    expect(res.statusCode).toBe(200);
    expect(res._getJSONData()).toEqual(
      expect.objectContaining({
        id: testUserItem.ID,
        username: testUserItem.Username,
        name: testUserItem.Name,
      }),
    );
  });
});
