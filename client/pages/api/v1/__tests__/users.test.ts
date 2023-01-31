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

jest.mock('ddb/users');
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
      EntityType: 'user' as const,
      ID: 'test-user',
      Username: 'testUser',
      Name: 'Test User',
      [gsi1PartitionKey]: 'users' as const,
      [gsi1SortKey]: 'user#test-user',
    };
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
    expect(res._getJSONData()).toEqual([testUserItem]);
  });
});
