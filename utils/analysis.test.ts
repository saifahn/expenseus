import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { Transaction } from 'types/Transaction';
import {
  personalTotalsByMainCategory,
  personalTotalsBySubcategory,
} from './analysis';

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

  test('it returns the expected results based on both txns and shared txns', () => {
    const testTxns: (SharedTxn | Transaction)[] = [
      {
        id: 'test-shared-txn',
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
      {
        id: 'test-txn',
        userId: 'test-user',
        location: 'burger king',
        amount: 980,
        date: 123456,
        category: 'food.eating-out',
        details: '',
      },
      {
        id: 'test-clothing-txn',
        userId: 'test-user',
        location: 'BEAMS',
        amount: 5000,
        date: 234567,
        category: 'clothing.clothing',
        details: '',
      },
    ];
    const totals = personalTotalsByMainCategory('test-user', testTxns);
    expect(totals).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ category: 'clothing', total: 10400 }),
        expect.objectContaining({ category: 'food', total: 980 }),
      ]),
    );
  });
});

describe('personalTotalsBySubcategory', () => {
  test('it returns the an empty object when given no transactions', () => {
    const totals = personalTotalsBySubcategory('test-user', []);
    expect(totals).toEqual([]);
  });

  test('it returns the expected results for shared txn with split', () => {
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
    const totals = personalTotalsBySubcategory('test-user', testTxns);
    expect(totals).toContainEqual(
      expect.objectContaining({ category: 'clothing.clothing', total: 5400 }),
    );
  });
});
