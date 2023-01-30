import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { ItemDoesNotExistError } from './errors';
import { makeSharedTxnRepository, SharedTxnItem } from './sharedTxns';

const sharedTxnsTestTable = 'shared-txns-test-table';
const d = setUpDdb(sharedTxnsTestTable);
const {
  createSharedTxn,
  updateSharedTxn,
  deleteSharedTxn,
  getTxnsByTracker,
  getTxnsByTrackerBetweenDates,
  getTxnsByUserBetweenDates,
  getUnsettledTxnsByTracker,
  settleTxns,
} = makeSharedTxnRepository(d);

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
    let txns = await getTxnsByTracker('test-tracker');
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
    await createSharedTxn(initialTxnDetails);

    txns = await getTxnsByTracker('test-tracker');
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
    await createSharedTxn(initialTxnDetails);
    let txns = await getTxnsByTracker('test-tracker');
    const createdTxn = txns[0];

    const updatedTxnDetails = {
      ...initialTxnDetails,
      id: createdTxn.ID,
      location: 'somewhere-else',
      amount: 99999,
    };
    await updateSharedTxn(updatedTxnDetails);
    txns = await getTxnsByTracker('test-tracker');
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
    expect(updateSharedTxn(txnDetails)).rejects.toThrow(ItemDoesNotExistError);
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
    await createSharedTxn(initialTxnDetails);
    let txns = await getTxnsByTracker('test-tracker');
    expect(txns).toHaveLength(1);
    const createdTxn = txns[0];

    await deleteSharedTxn({
      tracker: createdTxn.Tracker,
      txnId: createdTxn.ID,
      participants: createdTxn.Participants,
    });
    txns = await getTxnsByTracker('test-tracker');
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
    await createSharedTxn(unsettledTxn);
    await createSharedTxn(settledTxn);

    const txns = await getUnsettledTxnsByTracker('test-tracker');
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
    await createSharedTxn(initialTxn);
    let txns = await getUnsettledTxnsByTracker('test-tracker');
    expect(txns).toHaveLength(0);
    txns = await getTxnsByTracker('test-tracker');

    let updatedTxn: SharedTxn = {
      ...initialTxn,
      id: txns[0].ID,
      unsettled: true,
    };
    await updateSharedTxn(updatedTxn);
    txns = await getUnsettledTxnsByTracker('test-tracker');
    expect(txns).toHaveLength(1);

    // mark it as settled manually
    updatedTxn = {
      ...updatedTxn,
      unsettled: false,
    };
    await updateSharedTxn(updatedTxn);
    txns = await getUnsettledTxnsByTracker('test-tracker');
    expect(txns).toHaveLength(0);
  });

  test('all shared txns can be settled at once', async () => {
    // create two unsettled transactions
    const first: SharedTxn = {
      date: 1000 * 1000,
      location: '',
      amount: 34567,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
      unsettled: true,
    };
    const second: SharedTxn = {
      date: 2000 * 1000,
      location: '',
      amount: 99999,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
      unsettled: true,
    };
    await createSharedTxn(first);
    await createSharedTxn(second);
    // check the length of the unsettled txns
    let txns = await getUnsettledTxnsByTracker('test-tracker');
    expect(txns).toHaveLength(2);
    // trigger the settling
    const settleInput = txns.map((txn) => ({
      id: txn.ID,
      trackerId: txn.Tracker,
      participants: txn.Participants,
    }));
    await settleTxns(settleInput);
    // check the length of unsettled txns
    txns = await getUnsettledTxnsByTracker('test-tracker');
    expect(txns).toHaveLength(0);
  });

  test('txns are retrieved correctly by tracker between dates', async () => {
    const first: SharedTxn = {
      date: 1000 * 1000,
      location: '',
      amount: 34567,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    const second: SharedTxn = {
      date: 2000 * 1000,
      location: '',
      amount: 99999,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'a-different-tracker',
      details: '',
    };
    await createSharedTxn(first);
    await createSharedTxn(second);

    let input = {
      tracker: first.tracker,
      from: 1000 * 1000,
      to: 1500 * 1000,
    };
    let txns = await getTxnsByTrackerBetweenDates(input);
    expect(txns).toHaveLength(1);
    assertContainsTxnWithEqualDetails(txns, first);

    // same tracker, no txns in date range
    input = {
      tracker: first.tracker,
      from: 2000 * 1000,
      to: 3000 * 1000,
    };
    txns = await getTxnsByTrackerBetweenDates(input);
    expect(txns).toHaveLength(0);

    // different tracker, no results in date range
    input = {
      tracker: 'no-results-tracker',
      from: 1000 * 1000,
      to: 1500 * 1000,
    };
    txns = await getTxnsByTrackerBetweenDates(input);
    expect(txns).toHaveLength(0);
  });

  test('txns are retrieved correctly by user between dates', async () => {
    const first: SharedTxn = {
      date: 1000 * 1000,
      location: '',
      amount: 34567,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-02'],
      tracker: 'test-tracker',
      details: '',
    };
    const second: SharedTxn = {
      date: 2000 * 1000,
      location: '',
      amount: 99999,
      category: 'unspecified.unspecified',
      payer: 'user-01',
      participants: ['user-01', 'user-03'],
      tracker: 'a-different-tracker',
      details: '',
    };
    await createSharedTxn(first);
    await createSharedTxn(second);

    let input = {
      user: 'user-01',
      from: 1000 * 1000,
      to: 2500 * 1000,
    };
    let txns = await getTxnsByUserBetweenDates(input);
    expect(txns).toHaveLength(2);
    assertContainsTxnWithEqualDetails(txns, first);

    // same user, outside dates
    input = {
      user: 'user-01',
      from: 3000 * 1000,
      to: 4000 * 1000,
    };
    txns = await getTxnsByUserBetweenDates(input);
    expect(txns).toHaveLength(0);

    // different user, no results
    input = {
      user: 'no-results-user',
      from: 1000 * 1000,
      to: 1500 * 1000,
    };
    txns = await getTxnsByUserBetweenDates(input);
    expect(txns).toHaveLength(0);
  });
});
