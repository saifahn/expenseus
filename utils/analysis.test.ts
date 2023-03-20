import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { Transaction } from 'types/Transaction';
import { personalTotalsByMainCategory } from './analysis';

describe('personalTotalsByMainCategory', () => {
  test('it returns the an empty object when given no transactions', () => {
    const totals = personalTotalsByMainCategory('test-user', []);
    expect(totals).toEqual([]);
  });

  test('it returns the expected results based on txns with splits', () => {
    const testTxns: (SharedTxn | Transaction)[] = [
      {
        id: 'test-txn',
        date: 123456,
        amount: 9000,
        location: 'clothes',
        tracker: 'test-tracker',
        category: 'clothing.clothing',
        participants: ['test-user', 'test-user-2'],
        details: '',
        payer: 'test-user',
        split: {
          'test-user': 0.6,
          'test-user-2': 0.4,
        },
      },
    ];
    const totals = personalTotalsByMainCategory('test-user', testTxns);
    expect(totals).toContainEqual(
      expect.objectContaining({ category: 'clothing', total: 5400 }),
    );
  });
});
