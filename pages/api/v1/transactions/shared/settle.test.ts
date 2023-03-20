import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import * as nextAuth from 'next-auth';
import { mockReqRes } from 'tests/api/common';
import { sharedTxnRepoFnsMock } from 'tests/api/doubles';
import { stringify } from 'uuid';
import settleTxnsHandler from './settle';

const nextAuthMock = jest.mocked(nextAuth);
jest.mock('ddb/sharedTxns');
const sharedTxnRepo = jest.mocked(makeSharedTxnRepository);

describe('settleTxnsHandler', () => {
  test('returns a 405 if called with a non-POST method', async () => {
    const { req, res } = mockReqRes('GET');
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('returns a 400 if the input is invalid', async () => {
    const { req, res } = mockReqRes('POST');
    req._setBody({ something: 'invalid' });
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(400);
  });

  test('returns a 401 if there is no valid session', async () => {
    const { req, res } = mockReqRes('POST');
    req.body = JSON.stringify([]);
    nextAuthMock.getServerSession.mockResolvedValueOnce(null);
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(401);
  });

  test('returns a 403 if there are any txns that the user does not belong to', async () => {
    const { req, res } = mockReqRes('POST');
    req.body = JSON.stringify([
      {
        id: 'test-txn',
        tracker: 'test-tracker',
        participants: ['test-user', 'test-user-2'],
      },
    ]);
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: { email: 'different-user' },
    });
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(403);
  });

  test('will successfully call the ddb function to settle txns', async () => {
    const { req, res } = mockReqRes('POST');
    const testSettleInput = [
      {
        id: 'test-txn',
        tracker: 'test-tracker',
        participants: ['test-user', 'test-user-2'],
      },
    ];
    req.body = JSON.stringify(testSettleInput);
    nextAuthMock.getServerSession.mockResolvedValueOnce({
      user: { email: 'test-user' },
    });
    sharedTxnRepo.mockReturnValueOnce(sharedTxnRepoFnsMock);
    await settleTxnsHandler(req, res);

    expect(res.statusCode).toBe(202);
    expect(sharedTxnRepoFnsMock.settleTxns).toHaveBeenCalledWith(
      testSettleInput,
    );
  });
});
