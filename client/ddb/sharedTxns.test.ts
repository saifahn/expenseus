import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { ItemDoesNotExistError } from './errors';
import {
  createSharedTxn,
  deleteSharedTxn,
  getTxnsByTracker,
  getUnsettledTxnsByTracker,
  SharedTxnItem,
  updateSharedTxn,
} from './sharedTxns';

const sharedTxnsTestTable = 'shared-txns-test-table';
const d = setUpDdb(sharedTxnsTestTable);

/**
 * A helper function to check that the retrieved txns contain a txn with the
 * same details. We can't compare them directly because the ID is missing in the
 * original txn, and the SharedTxnItem and SharedTxn have different properties
 */
function assertContainsTxnWithEqualDetails(
  txns: SharedTxnItem[],
  txn: SharedTxn,
) {
  expect(txns).toContainEqual(
    expect.objectContaining({
      Date: txn.date,
      Location: txn.location,
      Amount: txn.amount,
      Category: txn.category,
      Payer: txn.payer,
      Participants: txn.participants,
      Tracker: txn.tracker,
      Details: txn.details,
    }),
  );
}

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

    const initialTxnDetails: SharedTxn = {
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, initialTxnDetails);

    txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(1);
    assertContainsTxnWithEqualDetails(txns, initialTxnDetails);
  });

  test('a shared txn can be updated successfully', async () => {
    const initialTxnDetails: SharedTxn = {
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified',
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
    assertContainsTxnWithEqualDetails(txns, updatedTxnDetails);
  });

  test('an error will be thrown when trying to update a non-existent shared txn', async () => {
    const txnDetails: SharedTxn = {
      id: 'non-existent',
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    expect(updateSharedTxn(d, txnDetails)).rejects.toThrow(
      ItemDoesNotExistError,
    );
  });

  test('a shared txn can be deleted successfully', async () => {
    const initialTxnDetails: SharedTxn = {
      date: 1000 * 1000,
      location: 'somewhere',
      amount: 12345,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, initialTxnDetails);
    let txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(1);
    const createdTxn = txns[0];

    await deleteSharedTxn(d, {
      tracker: createdTxn.Tracker,
      txnId: createdTxn.ID,
      participants: createdTxn.Participants,
    });
    txns = await getTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(0);
  });

  test('it is possible to retrieve only unsettled txns', async () => {
    const unsettledTxn: SharedTxn = {
      date: 1000 * 1000,
      location: 'somewhere unsettling',
      amount: 12345,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
      unsettled: true,
    };
    const settledTxn: SharedTxn = {
      date: 1000 * 1000,
      location: 'this is the sound of already settling',
      amount: 345,
      category: 'beauty.cosmetics',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, unsettledTxn);
    await createSharedTxn(d, settledTxn);

    const txns = await getUnsettledTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(1);
    assertContainsTxnWithEqualDetails(txns, unsettledTxn);
  });

  test('a shared txn can be updated to be marked as settled successfully', async () => {
    const initialTxn: SharedTxn = {
      date: 1000 * 1000,
      location: '',
      amount: 34567,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    await createSharedTxn(d, initialTxn);
    let txns = await getUnsettledTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(0);
    txns = await getTxnsByTracker(d, 'test-tracker');

    let updatedTxn: SharedTxn = {
      ...initialTxn,
      id: txns[0].ID,
      unsettled: true,
    };
    await updateSharedTxn(d, updatedTxn);
    txns = await getUnsettledTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(1);

    // mark it as settled manually
    updatedTxn = {
      ...updatedTxn,
      unsettled: false,
    };
    await updateSharedTxn(d, updatedTxn);
    txns = await getUnsettledTxnsByTracker(d, 'test-tracker');
    expect(txns).toHaveLength(0);
  });
});
