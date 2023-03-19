import {
  MainCategoryKey,
  subcategories,
  SubcategoryKey,
} from 'data/categories';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { Transaction } from 'types/Transaction';

export function calculateTotal(txns: (Transaction | SharedTxn)[]) {
  let total = 0;
  for (const txn of txns) {
    total += txn.amount;
  }
  return total;
}

export function totalsByMainCategory(txns: (Transaction | SharedTxn)[]) {
  const totals = {} as Record<MainCategoryKey, number>;
  for (const txn of txns) {
    const mainCategory = subcategories[txn.category].mainCategory;
    if (!totals[mainCategory]) totals[mainCategory] = 0;
    totals[mainCategory] += txn.amount;
  }
  return Object.entries(totals).map(([category, total]) => ({
    category,
    total,
  })) as { category: MainCategoryKey; total: number }[];
}

export function totalsBySubCategory(txns: (Transaction | SharedTxn)[]) {
  const totals = {} as Record<SubcategoryKey, number>;
  for (const txn of txns) {
    if (!totals[txn.category]) totals[txn.category] = 0;
    totals[txn.category] += txn.amount;
  }
  return Object.entries(totals).map(([category, total]) => ({
    category,
    total,
  })) as { category: SubcategoryKey; total: number }[];
}

export function calculatePersonalTotal(
  user: string,
  txns: (Transaction | SharedTxn)[],
) {
  let total = 0;
  for (const txn of txns) {
    // if shared txn
    if ('split' in txn) {
      if (txn.split?.[user]) {
        total += txn.amount * txn.split[user];
        continue;
      }
      // if there's no split, use the default split
      total += txn.amount * (1 / txn.participants.length);
      continue;
    }
    total += txn.amount;
  }
  return total;
}
