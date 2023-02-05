import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';
import { makeSharedTxnRepository } from 'ddb/sharedTxns';
import { makeTrackerRepository } from 'ddb/trackers';

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

export const sharedTxnRepoFnsMock: ReturnType<typeof makeSharedTxnRepository> =
  {
    createSharedTxn: jest.fn(),
    updateSharedTxn: jest.fn(),
    deleteSharedTxn: jest.fn(),
    getTxnsByTracker: jest.fn(),
    getTxnsByTrackerBetweenDates: jest.fn(),
    getSharedTxnsByUserBetweenDates: jest.fn(),
    getUnsettledTxnsByTracker: jest.fn(),
    settleTxns: jest.fn(),
  };

export const trackerRepoFnsMock: ReturnType<typeof makeTrackerRepository> = {
  createTracker: jest.fn(),
  getTracker: jest.fn(),
  getTrackersByUser: jest.fn(),
};
