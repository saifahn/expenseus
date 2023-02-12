import { mockReqRes } from 'tests/api/common';
import getUnsettledTxnsByTrackerHandler from './unsettled';

describe('getUnsettledTxnsByTrackerHandler', () => {
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

  test.todo('it returns 404 when there is no tracker');
  test.todo('it returns 403 if the user is not part of the tracker');
  test.todo(
    'it returns the transactions, debtee as the logged in user, debtor as the other user in the tracker, and amount owed as a number',
  );
});
