import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';
import {
  makeSharedTxnRepository,
  SharedTxn,
  SharedTxnItem,
} from 'ddb/sharedTxns';
import { getServerSession } from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { sharedTxnRepoFnsMock } from 'tests/api/doubles';
import txnsByTrackerHandler, { CreateSharedTxnPayload } from './transactions';

jest.mock('ddb/sharedTxns');
const sharedTxnRepo = jest.mocked(makeSharedTxnRepository);
const sessionMock = jest.mocked(getServerSession);

function assertEqualSharedTxnDetails(txn: SharedTxn, item: SharedTxnItem) {
  expect(txn).toEqual(
    expect.objectContaining({
      id: item.ID,
      tracker: item.Tracker,
      date: item.Date,
      amount: item.Amount,
      participants: item.Participants,
      location: item.Location,
      category: item.Category,
      payer: item.Payer,
      details: item.Details,
      ...(item.Unsettled && { unsettled: true }),
    }),
  );
}

describe('txnsByTrackerHandler', () => {
  test('returns a 405 if called with a non-POST or GET method', async () => {
    const { req, res } = mockReqRes('DELETE');
    await txnsByTrackerHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('returns a 401 if called with no valid session', async () => {
    const { req, res } = mockReqRes('GET');
    await txnsByTrackerHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  describe('GET - get all txns from tracker', () => {
    test('it return all txns returned from the store for the tracker', async () => {
      const { req, res } = mockReqRes('GET');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const mockItem: SharedTxnItem = {
        [tablePartitionKey]: 'tracker#test-tracker',
        [tableSortKey]: 'txn.shared#test-shared-txn',
        [gsi1PartitionKey]: 'tracker#test-tracker',
        [gsi1SortKey]: 'txn.shared#12345678#test-shared-txn',
        EntityType: 'sharedTransaction',
        ID: 'test-shared-txn',
        Tracker: 'test-tracker',
        Date: 12345678,
        Amount: 2473,
        Participants: ['test-user', 'test-user-2'],
        Location: 'LIFE',
        Category: 'food.groceries',
        Payer: 'test-user',
        Details: '',
      };
      const getSharedTxnMock = jest.fn().mockResolvedValueOnce([mockItem]);
      sharedTxnRepo.mockReturnValueOnce({
        ...sharedTxnRepoFnsMock,
        getTxnsByTracker: getSharedTxnMock,
      });
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(200);
      const result = res._getJSONData();
      expect(result).toHaveLength(1);
      assertEqualSharedTxnDetails(result[0], mockItem);
    });

    test.todo('it returns a 404 if the tracker does not exist');
  });

  describe('POST - create shared txn', () => {
    test('it returns a 400 with an invalid input', async () => {
      const { req, res } = mockReqRes('POST');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      // mock the input
      const createSharedTxnInput = {
        invalid: 'input',
      };
      req._setBody(createSharedTxnInput);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(400);
    });
    test('it returns a 403 when trying to create a shared txn without the session user as a participant', async () => {
      const { req, res } = mockReqRes('POST');
      sessionMock.mockResolvedValueOnce({ user: { email: 'different-user' } });
      const createSharedTxnInput: CreateSharedTxnPayload = {
        date: 12345678,
        amount: 2473,
        location: 'LIFE',
        category: 'food.groceries',
        tracker: 'test-tracker',
        participants: ['test-user', 'test-user-2'],
        payer: 'test-user',
        details: '',
      };
      req._setBody(createSharedTxnInput);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(403);
    });

    test('it successfully creates a shared txn', async () => {
      const { req, res } = mockReqRes('POST');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const createSharedTxnInput: CreateSharedTxnPayload = {
        date: 12345678,
        amount: 2473,
        location: 'LIFE',
        category: 'food.groceries',
        tracker: 'test-tracker',
        participants: ['test-user', 'test-user-2'],
        payer: 'test-user',
        details: '',
      };
      req._setBody(createSharedTxnInput);
      sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(sharedTxnRepoFnsMock.createSharedTxn).toHaveBeenCalledWith(
        createSharedTxnInput,
      );
    });
  });
});
