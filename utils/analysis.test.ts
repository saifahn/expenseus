import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { Transaction } from 'types/Transaction';
import {
  personalTotalsByMainCategory,
  personalTotalsBySubcategory,
  totalsByCategory,
} from './analysis';

describe('totalsByCategory', () => {
  test('it returns the expected results based on txns', () => {
    const testTxns: Transaction[] = [
      {
        id: 'test-txn',
        date: 123456,
        amount: 9000,
        location: 'clothes',
        category: 'clothing.clothing',
        details: '',
        userId: 'test-user-1',
      },
      {
        id: 'some-shoes',
        date: 234567,
        amount: 12300,
        location: 'allbirds',
        category: 'clothing.footwear',
        details: '',
        userId: 'test-user-1',
      },
      {
        id: 'groceries',
        date: 123456,
        amount: 3281,
        location: 'mybasket',
        category: 'food.groceries',
        userId: 'test-user-1',
        details: '',
      },
      {
        id: 'fancy-restaurant',
        date: 329142,
        amount: 55280,
        location: 'jiro sushi',
        category: 'food.eating-out',
        userId: 'test-user-1',
        details: '',
      },
      {
        id: 'conbini drink',
        date: 123456,
        amount: 230,
        location: 'seven eleven',
        category: 'food.food',
        userId: 'test-user-1',
        details: '',
      },
    ];
    const totals = totalsByCategory(testTxns);
    expect(totals).toContainEqual({
      mainCategory: 'clothing',
      total: 21_300,
      subcategories: expect.arrayContaining([
        {
          category: 'clothing.footwear',
          total: 12_300,
        },
        { category: 'clothing.clothing', total: 9000 },
      ]),
    });
    expect(totals).toContainEqual({
      mainCategory: 'food',
      total: 58_791,
      subcategories: expect.arrayContaining([
        {
          category: 'food.food',
          total: 230,
        },
        { category: 'food.eating-out', total: 55_280 },
        { category: 'food.groceries', total: 3281 },
      ]),
    });
  });
});

describe('personalTotalsByMainCategory', () => {
  test('it returns an empty array when given no transactions', () => {
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

  test('it returns the expected results for txns and shared txns', () => {
    const testTxns: (SharedTxn | Transaction)[] = [
      {
        id: 'test-shared-txn',
        date: 123456,
        amount: 9000,
        location: 'adidas',
        tracker: 'test-tracker',
        category: 'clothing.footwear',
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
    const totals = personalTotalsBySubcategory('test-user', testTxns);
    expect(totals).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ category: 'clothing.clothing', total: 5000 }),
        expect.objectContaining({ category: 'clothing.footwear', total: 5400 }),
        expect.objectContaining({ category: 'food.eating-out', total: 980 }),
      ]),
    );
  });
});
