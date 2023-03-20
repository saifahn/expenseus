import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';
import { makeTrackerRepository } from 'ddb/trackers';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { trackerRepoFnsMock } from 'tests/api/doubles';
import getTrackersByUserHandler from './[userId]';

jest.mock('ddb/trackers');
const trackerRepo = jest.mocked(makeTrackerRepository);
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

  test('it returns a list of trackers for a user successfully', async () => {
    const { req, res } = mockReqRes('GET');
    sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    req.query = { userId: 'test-user' };
    const getTrackerMock: ReturnType<typeof trackerRepo>['getTrackersByUser'] =
      jest.fn(async () => [
        {
          [tablePartitionKey]: 'tracker#test-tracker',
          [tableSortKey]: 'tracker#test-tracker',
          [gsi1PartitionKey]: 'trackers' as const,
          [gsi1SortKey]: 'tracker#test-tracker',
          EntityType: 'tracker' as const,
          ID: 'test-tracker',
          Name: 'Test Tracker',
          Users: ['test-user', 'test-user-2'],
        },
      ]);
    trackerRepo.mockReturnValueOnce({
      ...trackerRepoFnsMock,
      getTrackersByUser: getTrackerMock,
    });
    await getTrackersByUserHandler(req, res);

    expect(res.statusCode).toBe(200);
    expect(res._getJSONData()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          id: 'test-tracker',
          name: 'Test Tracker',
          users: ['test-user', 'test-user-2'],
        }),
      ]),
    );
  });
});
