import { NextApiRequest, NextApiResponse } from 'next';
import { createMocks, RequestMethod } from 'node-mocks-http';
import usersHandler from '../users';
import { makeUserRepository } from 'ddb/users';
import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';

jest.mock('ddb/users', () => {
  const original = jest.requireActual('ddb/users');
  return {
    ...original,
    makeUserRepository: jest.fn(),
  };
});

const usersRepo = jest.mocked(makeUserRepository);

describe('/api/v1/users API endpoint', () => {
  function mockReqRes(method: RequestMethod = 'GET') {
    const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
      method,
    });
    return { req, res };
  }

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
        getAllUsers: jest.fn(async () => [testUserItem]),
      };
    });

    const { req, res } = mockReqRes();
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

  it('returns an error when called with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await usersHandler(req, res);

    expect(res.statusCode).toBe(405);
  });
});
