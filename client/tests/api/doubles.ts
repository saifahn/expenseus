import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';

export const mockTxnItem = {
  [tablePartitionKey]: 'user#test-user',
  [tableSortKey]: 'txn#test-txn',
  [gsi1PartitionKey]: 'user#test-user',
  [gsi1SortKey]: 'txn#12345678#test-txn',
  EntityType: 'transaction' as const,
  ID: 'test-txn',
  UserID: 'test-user',
  Date: 12345678,
  Amount: 9275,
  Location: 'somewhere',
  Category: 'unspecified.unspecified' as const,
  Details: '',
};
