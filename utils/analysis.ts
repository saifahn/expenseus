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

type TotalsGroupedByCategoryDict = {
  [k in MainCategoryKey]: {
    [x in SubcategoryKey]: number;
  };
};

type TotalByCategory = {
  mainCategory: MainCategoryKey;
  total: number;
  subcategories: Array<{
    category: SubcategoryKey;
    total: number;
  }>;
};

/**
 * Converts a list of transactions or shared transactions into a data structure
 * to make it easier to display values grouped by main category on the analysis
 * pages
 */
export function totalsByCategory(txns: Transaction[] | SharedTxn[]) {
  // create a dictionary first as it is easier to manage the unknown categories
  const totalsDict = {} as TotalsGroupedByCategoryDict;
  for (const txn of txns) {
    const mainCategory = subcategories[txn.category].mainCategory;
    if (!totalsDict[mainCategory]) {
      totalsDict[mainCategory] =
        {} as TotalsGroupedByCategoryDict[MainCategoryKey];
    }
    if (!totalsDict[mainCategory][txn.category]) {
      totalsDict[mainCategory][txn.category] = 0;
    }
    totalsDict[mainCategory][txn.category] += txn.amount;
  }

  // convert to an array because it will be easier to use on the FE
  const totals = (Object.keys(totalsDict) as MainCategoryKey[]).map(
    (mainCatKey) => {
      // for each main category, we will create a total for that category
      // and its subcategories
      const subcatKeys = Object.keys(
        totalsDict[mainCatKey],
      ) as SubcategoryKey[];

      const categoryTotal: TotalByCategory = {
        mainCategory: mainCatKey,
        total: 0,
        subcategories: [],
      };

      for (const subcatKey of subcatKeys) {
        const subcategoryTotal = totalsDict[mainCatKey][subcatKey];
        categoryTotal.subcategories.push({
          category: subcatKey,
          total: subcategoryTotal,
        });
        categoryTotal.total += subcategoryTotal;
      }

      return categoryTotal;
    },
  );

  return totals;
}

export function calculatePersonalTotal(
  user: string,
  txns: (Transaction | SharedTxn)[],
) {
  let total = 0;

  for (const txn of txns) {
    const amount = calculateAmountOwedByUser(txn, user);
    total += amount;
  }

  return total;
}

export function personalTotalsByCategory(
  user: string,
  txns: (Transaction | SharedTxn)[],
): TotalByCategory[] {
  // create a dictionary first as it is easier to manage the unknown categories
  const totalsDict = {} as TotalsGroupedByCategoryDict;
  for (const txn of txns) {
    const mainCategory = subcategories[txn.category].mainCategory;
    if (!totalsDict[mainCategory]) {
      totalsDict[mainCategory] =
        {} as TotalsGroupedByCategoryDict[MainCategoryKey];
    }
    if (!totalsDict[mainCategory][txn.category]) {
      totalsDict[mainCategory][txn.category] = 0;
    }
    const amount = calculateAmountOwedByUser(txn, user);
    totalsDict[mainCategory][txn.category] += amount;
  }

  // convert to an array because it will be easier to use on the FE
  const totals = (Object.keys(totalsDict) as MainCategoryKey[]).map(
    (mainCatKey) => {
      // for each main category, we will create a total for that category
      // and its subcategories
      const subcatKeys = Object.keys(
        totalsDict[mainCatKey],
      ) as SubcategoryKey[];

      const categoryTotal: TotalByCategory = {
        mainCategory: mainCatKey,
        total: 0,
        subcategories: [],
      };

      for (const subcatKey of subcatKeys) {
        const subcategoryTotal = totalsDict[mainCatKey][subcatKey];
        categoryTotal.subcategories.push({
          category: subcatKey,
          total: subcategoryTotal,
        });
        categoryTotal.total += subcategoryTotal;
      }

      return categoryTotal;
    },
  );

  return totals;
}

/**
 * Helper function that takes a user and a txn and returns the amount
 * that user is responsible for for that txn.
 */
function calculateAmountOwedByUser(txn: SharedTxn | Transaction, user: string) {
  // if sharedTxn
  if ('split' in txn) {
    const split = txn.split?.[user];
    return split
      ? txn.amount * split
      : txn.amount * (1 / txn.participants.length);
  } else {
    // if regular txn
    return txn.amount;
  }
}
