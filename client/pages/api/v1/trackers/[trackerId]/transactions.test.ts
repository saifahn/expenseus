import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import { getServerSession } from 'next-auth';
import { assertEqualSharedTxnDetails, mockReqRes } from 'tests/api/common';
import { mockSharedTxnItem, sharedTxnRepoFnsMock } from 'tests/api/doubles';
import txnsByTrackerHandler from './transactions';

jest.mock('ddb/sharedTxns');
const sharedTxnRepo = jest.mocked(makeSharedTxnRepository);
const sessionMock = jest.mocked(getServerSession);

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
      const getSharedTxnMock = jest
        .fn()
        .mockResolvedValueOnce([mockSharedTxnItem]);
      sharedTxnRepo.mockReturnValueOnce({
        ...sharedTxnRepoFnsMock,
        getTxnsByTracker: getSharedTxnMock,
      });
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(200);
      const result = res._getJSONData();
      expect(result).toHaveLength(1);
      assertEqualSharedTxnDetails(result[0], mockSharedTxnItem);
    });

    test('it returns txns with split values correctly', async () => {
      const { req, res } = mockReqRes('GET');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const getSharedTxnMock = jest
        .fn()
        .mockResolvedValueOnce([
          {
            ...mockSharedTxnItem,
            SplitJSON: '{"test-user":0.6, "test-user-2":0.4}',
          },
        ]);
      sharedTxnRepo.mockReturnValueOnce({
        ...sharedTxnRepoFnsMock,
        getTxnsByTracker: getSharedTxnMock,
      });
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(200);
      const result = res._getJSONData();
      expect(result).toHaveLength(1);
      assertEqualSharedTxnDetails(result[0], mockSharedTxnItem);
      expect(result[0]).toEqual(
        expect.objectContaining({
          split: {
            'test-user': 0.6,
            'test-user-2': 0.4,
          },
        }),
      );
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
      req.body = JSON.stringify(createSharedTxnInput);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(400);
    });
    test('it returns a 403 when trying to create a shared txn without the session user as a participant', async () => {
      const { req, res } = mockReqRes('POST');
      sessionMock.mockResolvedValueOnce({ user: { email: 'different-user' } });
      const createSharedTxnInput = {
        date: 12345678,
        amount: 2473,
        location: 'LIFE',
        category: 'food.groceries',
        participants: ['test-user', 'test-user-2'],
        payer: 'test-user',
        details: '',
      };
      req.query.trackerId = 'test-tracker';
      req.body = JSON.stringify(createSharedTxnInput);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(403);
    });

    test('it successfully creates a shared txn', async () => {
      const { req, res } = mockReqRes('POST');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const createSharedTxnInput = {
        date: 12345678,
        amount: 2473,
        location: 'LIFE',
        category: 'food.groceries',
        participants: ['test-user', 'test-user-2'],
        payer: 'test-user',
        details: '',
      };
      req.query.trackerId = 'test-tracker';
      req.body = JSON.stringify(createSharedTxnInput);
      sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(sharedTxnRepoFnsMock.createSharedTxn).toHaveBeenCalledWith({
        ...createSharedTxnInput,
        tracker: 'test-tracker',
      });
    });

    test('it successfully creates a shared txn with a split', async () => {
      const { req, res } = mockReqRes('POST');
      sessionMock.mockResolvedValueOnce({ user: { email: 'test-user' } });
      const createSharedTxnInput = {
        date: 12345678,
        amount: 2473,
        location: 'LIFE',
        category: 'food.groceries',
        participants: ['test-user', 'test-user-2'],
        payer: 'test-user',
        details: '',
        split: {
          'test-user': 0.4,
          'test-user-2': 0.6,
        },
      };
      req.query.trackerId = 'test-tracker';
      req.body = JSON.stringify(createSharedTxnInput);
      sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
      await txnsByTrackerHandler(req, res);

      expect(res.statusCode).toBe(202);
      expect(sharedTxnRepoFnsMock.createSharedTxn).toHaveBeenCalledWith({
        ...createSharedTxnInput,
        tracker: 'test-tracker',
        split: {
          'test-user': 0.4,
          'test-user-2': 0.6,
        },
      });
    });
  });
});
