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
});
