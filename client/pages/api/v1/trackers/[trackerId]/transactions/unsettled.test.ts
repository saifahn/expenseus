import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { sharedTxnRepoFnsMock } from 'tests/api/doubles';
import getUnsettledTxnsByTrackerHandler from './unsettled';

jest.mock('ddb/sharedTxns');
const sharedTxnRepo = jest.mocked(makeSharedTxnRepository);
const sessionMock = jest.mocked(getServerSession);

describe('getUnsettledTxnsByTrackerHandler', () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  test('it returns 405 on a non-GET method', async () => {
    const { req, res } = mockReqRes('POST');
    await getUnsettledTxnsByTrackerHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('it returns 401 when there is no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    await getUnsettledTxnsByTrackerHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('it returns the transactions, debtee as the logged in user, debtor as the other user in the tracker, and amount owed as a number', async () => {
    const { req, res } = mockReqRes('GET');
    sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
    req.query = {
      trackerId: 'test-tracker',
    };
    sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
    await getUnsettledTxnsByTrackerHandler(req, res);

    expect(sharedTxnRepoFnsMock.getUnsettledTxnsByTracker).toHaveBeenCalledWith(
      'test-tracker',
    );
  });

  test.todo('it returns 403 if the user is not part of the tracker');
});
