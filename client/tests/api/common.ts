import { TxnItem } from 'ddb/txns';
import { NextApiRequest, NextApiResponse } from 'next';
import { RequestMethod, createMocks } from 'node-mocks-http';
import { Transaction } from 'types/Transaction';

export function mockReqRes(method: RequestMethod) {
  const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
    method,
  });
  return { req, res };
}

/**
 * helper function to assert details from a txn match txn item
 */
export function assertEqualTxnDetails(txn: Transaction, txnItem: TxnItem) {
  expect(txn).toEqual(
    expect.objectContaining({
      userId: txnItem.UserID,
      location: txnItem.Location,
      amount: txnItem.Amount,
      date: txnItem.Date,
      category: txnItem.Category,
      details: txnItem.Details,
    }),
  );
}
