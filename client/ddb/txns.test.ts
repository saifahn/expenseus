import { ItemDoesNotExistError } from 'ddb/errors';
import {
  setUpDdb,
  createTableIfNotExists,
  TESTONLY_deleteTable,
} from 'ddb/schema';
import { makeTxnRepository, TxnItem } from 'ddb/txns';
import { CreateTxnPayload } from 'pages/api/v1/transactions';

const txnTestTable = 'txn-test-table';
const d = setUpDdb(txnTestTable);
const {
  createTxn,
  deleteTxn,
  getBetweenDates,
  getTxn,
  getTxnsByUserId,
  updateTxn,
} = makeTxnRepository(d);

// helper function to assert details from txnItem match an original txn
function assertEqualDetails(txnItem: TxnItem, txn: CreateTxnPayload) {
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
    await TESTONLY_deleteTable(txnTestTable);
  });

  test('a txn can be created successfully', async () => {
    let txns = await getTxnsByUserId('test-user');
    expect(txns).toHaveLength(0);

    const testTxn = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000000,
      category: 'unspecified.unspecified',
      details: '',
    } as const;
    await createTxn(testTxn);

    txns = await getTxnsByUserId('test-user');
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
    await createTxn(testTxn);
    const txns = await getTxnsByUserId('test-user');
    const createdTxn = txns[0];

    const result = await getTxn({
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
    await createTxn(testTxn);
    let txns = await getTxnsByUserId('test-user');
    const createdTxn = txns[0];
    assertEqualDetails(createdTxn, testTxn);

    const updatedTxn = {
      ...testTxn,
      id: createdTxn.ID,
      location: 'updated-location',
      amount: 999999,
    };
    await updateTxn(updatedTxn);
    txns = await getTxnsByUserId('test-user');
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
    expect(updateTxn(updatedTxn)).rejects.toThrowError(ItemDoesNotExistError);
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
    await createTxn(testTxn);
    let txns = await getTxnsByUserId('test-user');
    expect(txns).toHaveLength(1);

    await deleteTxn({ txnId: txns[0].ID, userId: txns[0].UserID });

    txns = await getTxnsByUserId('test-user');
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
    await createTxn(testTxn);

    let txns = await getBetweenDates({
      userId: testTxn.userId,
      from: 1000 * 1000,
      to: 1000 * 1500,
    });
    expect(txns).toHaveLength(1);
    assertEqualDetails(txns[0], testTxn);

    // a date range outside returns none
    txns = await getBetweenDates({
      userId: testTxn.userId,
      from: 2000 * 1000,
      to: 2000 * 1500,
    });
    expect(txns).toHaveLength(0);
  });
});
