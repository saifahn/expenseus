import { Transaction } from 'types/Transaction';
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
