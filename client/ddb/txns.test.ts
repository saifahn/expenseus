import { ItemDoesNotExistError } from 'ddb/errors';
import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import {
  createTxn,
  deleteTxn,
  getBetweenDates,
  getTxn,
  getTxnsByUserId,
  TxnItem,
  updateTxn,
} from 'ddb/txns';
import { Transaction } from 'types/Transaction';

const txnTestTable = 'txn-test-table';
const d = setUpDdb(txnTestTable);

// helper function to assert details from txnItem match an original txn
function assertEqualDetails(txnItem: TxnItem, txn: Transaction) {
  expect(txnItem).toEqual(
    expect.objectContaining({
      UserID: txn.userId,
      Location: txn.location,
      Amount: txn.amount,
      Date: txn.date,
      Category: txn.category,
      Details: txn.details,
    }),
  );
}

describe('Transactions', () => {
  beforeEach(async () => {
    await createTableIfNotExists(txnTestTable);
  });

  afterEach(async () => {
    await deleteTable(txnTestTable);
  });

  test('a txn can be created successfully', async () => {
    let txns = await getTxnsByUserId(d, 'test-user');
    expect(txns).toHaveLength(0);

    const testTxn = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    await createTxn(d, testTxn);

    txns = await getTxnsByUserId(d, 'test-user');
    expect(txns).toHaveLength(1);
    assertEqualDetails(txns[0], testTxn);
  });

  test('a txn can be retrieved successfully', async () => {
    const testTxn = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    await createTxn(d, testTxn);
    const txns = await getTxnsByUserId(d, 'test-user');
    const createdTxn = txns[0];

    const result = await getTxn(d, {
      txnId: createdTxn.ID,
      userId: testTxn.userId,
    });

    assertEqualDetails(result, testTxn);
  });

  test('a txn can be updated successfully', async () => {
    const testTxn = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    await createTxn(d, testTxn);
    let txns = await getTxnsByUserId(d, 'test-user');
    const createdTxn = txns[0];
    assertEqualDetails(createdTxn, testTxn);

    const updatedTxn = {
      ...testTxn,
      id: createdTxn.ID,
      location: 'updated-location',
      amount: 999999,
    };
    await updateTxn(d, updatedTxn);
    txns = await getTxnsByUserId(d, 'test-user');
    assertEqualDetails(txns[0], updatedTxn);
  });

  test('a ItemDoesNotExist error will be thrown if a txn that does not exist is attempted to be updated', async () => {
    const updatedTxn = {
      id: 'non-existent-txn',
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    expect(updateTxn(d, updatedTxn)).rejects.toThrowError(
      ItemDoesNotExistError,
    );
  });

  test('a txn can be deleted successfully', async () => {
    const testTxn = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    await createTxn(d, testTxn);
    let txns = await getTxnsByUserId(d, 'test-user');
    expect(txns).toHaveLength(1);

    await deleteTxn(d, { txnId: txns[0].ID, userId: txns[0].UserID });

    txns = await getTxnsByUserId(d, 'test-user');
    expect(txns).toHaveLength(0);
  });

  test('txns can be retrieved for a given user between a date range', async () => {
    const testTxn = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    await createTxn(d, testTxn);

    let txns = await getBetweenDates(d, {
      userId: testTxn.userId,
      from: 1000 * 1000,
      to: 1000 * 1500,
    });
    expect(txns).toHaveLength(1);
    assertEqualDetails(txns[0], testTxn);

    // a date range outside returns none
    txns = await getBetweenDates(d, {
      userId: testTxn.userId,
      from: 2000 * 1000,
      to: 2000 * 1500,
    });
    expect(txns).toHaveLength(0);
  });
});