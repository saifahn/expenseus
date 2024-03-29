import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { Transaction } from 'types/Transaction';
import { ResponsiveBar } from '@nivo/bar';
import {
  categoryColours,
  MainCategoryKey,
  mainCategoryKeys,
  subcategories,
} from 'data/categories';
import { epochSecToUTCMonthEN, MonthEN } from 'utils/dates';
import { calculateAmountOwedByUser } from 'utils/analysis';

type TotalsForBarChart = Record<
  MonthEN,
  Partial<Record<MainCategoryKey, number>>
>;

/**
 * Takes a list of transactions and returns totals by month and main category
 * for use in data visualization.
 */
export function totalsForBarChart(txns: (Transaction | SharedTxn)[]) {
  const totals = {} as TotalsForBarChart;
  for (const txn of txns) {
    const month = epochSecToUTCMonthEN(txn.date);
    const mainCategory = subcategories[txn.category].mainCategory;
    if (!totals[month]) totals[month] = {};
    if (!totals[month][mainCategory]) totals[month][mainCategory] = 0;
    totals[month][mainCategory]! += txn.amount;
  }
  return totals;
}

export function personalTotalsForBarChart(
  txns: (Transaction | SharedTxn)[],
  user: string,
) {
  const totals = {} as TotalsForBarChart;
  for (const txn of txns) {
    const month = epochSecToUTCMonthEN(txn.date);
    const mainCategory = subcategories[txn.category].mainCategory;
    if (!totals[month]) totals[month] = {};
    if (!totals[month][mainCategory]) totals[month][mainCategory] = 0;
    const amount = calculateAmountOwedByUser(txn, user);
    totals[month][mainCategory]! += amount;
  }
  return totals;
}

export function BarChart(totals: TotalsForBarChart) {
  const entries = Object.entries(totals);
  const data = entries.map(([month, values]) => ({
    month,
    ...values,
  }));

  return (
    <ResponsiveBar
      keys={mainCategoryKeys}
      indexBy={'month'}
      data={data}
      margin={{
        top: 50,
        bottom: 50,
        right: 50,
        left: 50,
      }}
      colors={categoryColours}
      axisBottom={{
        legend: 'month',
        legendPosition: 'middle',
        legendOffset: 40,
      }}
      axisLeft={{
        format: (value: number) => value.toLocaleString(),
      }}
      valueFormat={(value) => value.toLocaleString()}
    />
  );
}
