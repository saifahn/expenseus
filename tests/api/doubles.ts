import {
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from 'ddb/schema';
import { makeSharedTxnRepository, SharedTxnItem } from 'ddb/sharedTxns';
import { makeTrackerRepository } from 'ddb/trackers';
import { makeUserRepository, UserItem } from 'ddb/users';

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

export const mockSharedTxnItem: SharedTxnItem = {
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

export const testUserItem: UserItem = {
  [tablePartitionKey]: 'user#test-user',
  [tableSortKey]: 'user#test-user',
  EntityType: 'user',
  ID: 'test-user',
  Username: 'testUser',
  Name: 'Test User',
  [gsi1PartitionKey]: 'users',
  [gsi1SortKey]: 'user#test-user',
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

export const userRepoFnsMock: ReturnType<typeof makeUserRepository> = {
  createUser: jest.fn(),
  getAllUsers: jest.fn(),
  getUser: jest.fn(),
};
