import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import {
  createSharedTxn,
  getTxnsByTracker,
  updateSharedTxn,
} from './sharedTxns';

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

    const initialTxnDetails = {
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified' as const,
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, initialTxnDetails);

    txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(1);
  });

  test('a shared txn can be updated successfully', async () => {
    const initialTxnDetails = {
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified' as const,
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, initialTxnDetails);
    let txns = await getTxnsByTracker(d, 'test-tracker');
    const createdTxn = txns[0];

    const updatedTxnDetails = {
      ...initialTxnDetails,
      id: createdTxn.ID,
      location: 'somewhere-else',
      amount: 99999,
    };
    await updateSharedTxn(d, updatedTxnDetails);
    txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toContainEqual(
      expect.objectContaining({
        Date: updatedTxnDetails.date,
        Location: updatedTxnDetails.location,
        Amount: updatedTxnDetails.amount,
        Category: updatedTxnDetails.category,
        Payer: updatedTxnDetails.payer,
        Participants: updatedTxnDetails.participants,
        Tracker: updatedTxnDetails.tracker,
        Details: updatedTxnDetails.details,
      }),
    );
  });
});
