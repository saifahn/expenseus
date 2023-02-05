import { Transaction } from 'types/Transaction';
import { SharedTxn, SharedTxnItem } from './sharedTxns';
import { TxnItem } from './txns';

export function txnItemToTxn(item: TxnItem): Transaction {
  return {
    id: item.ID,
    userId: item.UserID,
    location: item.Location,
    details: item.Details,
    amount: item.Amount,
    date: item.Date,
    category: item.Category,
  };
}

export function sharedTxnItemToModel(item: SharedTxnItem): SharedTxn {
  return {
    id: item.ID,
    tracker: item.Tracker,
    location: item.Location,
    details: item.Details,
    amount: item.Amount,
    date: item.Date,
    category: item.Category,
    participants: item.Participants,
    payer: item.Payer,
    ...(item.Unsettled && { unsettled: true }),
  };
}
