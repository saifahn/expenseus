import { makeSharedTxnRepository, SharedTxn } from 'ddb/sharedTxns';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { mockSharedTxnItem, sharedTxnRepoFnsMock } from 'tests/api/doubles';
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
    const getUnsettledMock: ReturnType<
      typeof sharedTxnRepo
    >['getUnsettledTxnsByTracker'] = jest.fn(async () => [
      {
        ...mockSharedTxnItem,
        Amount: 4000,
        Unsettled: 'X' as const,
      },
    ]);
    sharedTxnRepo.mockReturnValueOnce({
      ...sharedTxnRepoFnsMock,
      getUnsettledTxnsByTracker: getUnsettledMock,
    });
    await getUnsettledTxnsByTrackerHandler(req, res);

    expect(getUnsettledMock).toHaveBeenCalledWith('test-tracker');
    const result = res._getJSONData();
    expect(result).toEqual(
      expect.objectContaining({
        transactions: expect.arrayContaining([
          expect.objectContaining({
            tracker: mockSharedTxnItem.Tracker,
            date: mockSharedTxnItem.Date,
            participants: mockSharedTxnItem.Participants,
            amount: 4000,
            location: mockSharedTxnItem.Location,
            category: mockSharedTxnItem.Category,
            payer: mockSharedTxnItem.Payer,
          }),
        ]),
        debtor: 'test-user-2',
        debtee: 'test-user',
        amountOwed: 2000,
      }),
    );
  });

  test('it returns the right details when the logged in user is not the payer', async () => {
    const { req, res } = mockReqRes('GET');
    sessionMock.mockResolvedValueOnce({ user: { email: 'test-user-2' } });
    req.query = {
      trackerId: 'test-tracker',
    };
    const getUnsettledMock: ReturnType<
      typeof sharedTxnRepo
    >['getUnsettledTxnsByTracker'] = jest.fn(async () => [
      {
        ...mockSharedTxnItem,
        Amount: 4000,
        Unsettled: 'X' as const,
      },
    ]);
    sharedTxnRepo.mockReturnValueOnce({
      ...sharedTxnRepoFnsMock,
      getUnsettledTxnsByTracker: getUnsettledMock,
    });
    await getUnsettledTxnsByTrackerHandler(req, res);

    expect(getUnsettledMock).toHaveBeenCalledWith('test-tracker');
    const result = res._getJSONData();
    expect(result).toEqual(
      expect.objectContaining({
        debtor: 'test-user',
        debtee: 'test-user-2',
        amountOwed: -2000,
      }),
    );
  });

  test.todo('it returns 403 if the user is not part of the tracker');
});
