import usersHandler from './users';
import { makeUserRepository } from 'ddb/users';
import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';
import { mockReqRes } from 'tests/api/common';
import * as nextAuth from 'next-auth';

jest.mock('ddb/users', () => {
  const original = jest.requireActual('ddb/users');
  return {
    ...original,
    makeUserRepository: jest.fn(),
  };
});
const usersRepo = jest.mocked(makeUserRepository);

const nextAuthMock = jest.mocked(nextAuth);
const blankValidSession = {};

describe('/api/v1/users API endpoint', () => {
  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('should return a list of users when they are returned from the store', async () => {
    const testUserItem = {
      [tablePartitionKey]: 'user#test-user',
      [tableSortKey]: 'user#test-user',
      EntityType: 'user',
      ID: 'test-user',
      Username: 'testUser',
      Name: 'Test User',
      [gsi1PartitionKey]: 'users',
      [gsi1SortKey]: 'user#test-user',
    } as const;
    usersRepo.mockImplementationOnce(() => {
      return {
        createUser: jest.fn(),
        getUser: jest.fn(),
        getAllUsers: jest.fn(async () => [testUserItem]),
      };
    });

    nextAuthMock.getServerSession.mockResolvedValueOnce(blankValidSession);
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
    nextAuthMock.getServerSession.mockResolvedValueOnce(blankValidSession);
    await usersHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  it('returns a 401 error when there is no valid session', async () => {
    // mock the session to return invalid
    const { req, res } = mockReqRes('GET');
    nextAuthMock.getServerSession.mockResolvedValueOnce(null);
    await usersHandler(req, res);

    expect(res.statusCode).toBe(401);
  });
});
