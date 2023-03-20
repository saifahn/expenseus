import { personalTotalsByMainCategory } from './analysis';

describe('personalTotalsByMainCategory', () => {
  test('it returns the an empty object when given no transactions', () => {
    const totals = personalTotalsByMainCategory('test-user', []);
    expect(totals).toEqual([]);
  });
});
