import { makeTrackerRepository } from 'ddb/trackers';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { trackerRepoFnsMock } from 'tests/api/doubles';
import getTrackerByIdHandler from './[trackerId]';

jest.mock('ddb/trackers');
const trackerRepoMock = jest.mocked(makeTrackerRepository);
const serverSessionMock = jest.mocked(getServerSession);

describe('getTrackerByIdHandler', () => {
  test('it returns a 405 if called with a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getTrackerByIdHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns a 401 is no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    serverSessionMock.mockResolvedValueOnce(null);
    await getTrackerByIdHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('it returns a 404 if there is no tracker returned from the store', async () => {
    const { req, res } = mockReqRes('GET');
    serverSessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    trackerRepoMock.mockReturnValueOnce(trackerRepoFnsMock);
    await getTrackerByIdHandler(req, res);

    expect(res.statusCode).toBe(404);
  });

  test.todo(
    'it successfully returns a tracker when one is returned from the store',
  );
});
