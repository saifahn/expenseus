import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import { createSharedTxn, getTxnsByTracker } from './sharedTxns';

const sharedTxnsTestTable = 'shared-txns-test-table';
const d = setUpDdb(sharedTxnsTestTable);

describe('Shared Transactions', () => {
  beforeEach(async () => {
    await createTableIfNotExists(sharedTxnsTestTable);
  });

  afterEach(async () => {
    await deleteTable(sharedTxnsTestTable);
  });

  test('a shared txn can be created and retrieved correctly', async () => {
    let txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(0);

    const sharedTxnDetails = {
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified' as const,
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, sharedTxnDetails);

    txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(1);
  });
});
