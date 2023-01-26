import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import { createTxn, getTxnsByUserId, updateTxn } from 'ddb/txns';

const d = setUpDdb('test-table');

describe('Transactions', () => {
  beforeEach(async () => {
    await createTableIfNotExists('test-table');
  });

  afterEach(async () => {
    await deleteTable('test-table');
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
    expect(txns[0]).toEqual(
      expect.objectContaining({
        UserID: testTxn.userId,
        Location: testTxn.location,
        Amount: testTxn.amount,
        Date: testTxn.date,
        Category: testTxn.category,
        Details: testTxn.details,
      }),
    );
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
    expect(createdTxn).toEqual(
      expect.objectContaining({
        UserID: testTxn.userId,
        Location: testTxn.location,
        Amount: testTxn.amount,
        Date: testTxn.date,
        Category: testTxn.category,
        Details: testTxn.details,
      }),
    );

    const updatedTxn = {
      ...testTxn,
      id: createdTxn.ID,
      location: 'updated-location',
      amount: 999999,
    };
    await updateTxn(d, updatedTxn);
    txns = await getTxnsByUserId(d, 'test-user');
    expect(txns[0]).toEqual(
      expect.objectContaining({
        UserID: updatedTxn.userId,
        Location: updatedTxn.location,
        Amount: updatedTxn.amount,
        Date: updatedTxn.date,
        Category: updatedTxn.category,
        Details: updatedTxn.details,
      }),
    );
  });

  test.todo(
    'an error will be thrown if a txn that does not exist is attempted to be updated',
  );

  test.todo('a txn can be deleted successfully');
});
