import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import { createTxn, getTxnsByUserId } from 'ddb/txns';

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

    await createTxn(d, {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000000,
      category: 'unspecified.unspecified',
      details: '',
    });

    txns = await getTxnsByUserId(d, 'test-user');
    expect(txns).toHaveLength(1);
    // improvement: check that the cluck is the right one
  });
});
