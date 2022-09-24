import { useUserContext } from 'context/user';
import useSWR from 'swr';
import { Transaction } from 'types/Transaction';
import { Temporal } from 'temporal-polyfill';
import {
  epochSecToLocaleString,
  epochSecToUTCYear,
  plainDateStringToEpochSec,
} from 'utils/dates';
import { SharedTxn } from './shared/trackers/[trackerId]';
import { calculatePersonalTotal } from 'utils/analysis';
import { categoryNameFromKeyEN } from 'data/categories';

type AllTxnsResponse = {
  transactions: Transaction[];
  sharedTransactions: SharedTxn[];
};

export default function Home() {
  const { user } = useUserContext();
  const todayEpochSeconds = plainDateStringToEpochSec(
    Temporal.Now.plainDateISO().toString(),
  );
  const threeMonthsAgoEpochSeconds = plainDateStringToEpochSec(
    Temporal.Now.plainDateISO().subtract({ months: 3 }).toString(),
  );
  const { data: res, error } = useSWR<AllTxnsResponse>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}/all?from=${threeMonthsAgoEpochSeconds}&to=${todayEpochSeconds}`
      : null,
  );

  let txns: (Transaction | SharedTxn)[] = [];
  if (res) {
    txns = [...res.transactions, ...res.sharedTransactions].sort(
      (a, b) => b.date - a.date, // sorted descending, so the most recent dates are shown first
    );
  }
  const total = calculatePersonalTotal(user.id, txns);

  return (
    <>
      <section>
        <p className="mt-4">Hi, {user.username}!</p>
        {error && <p>Failed to load recent transactions</p>}
        {res === null && <p>Loading recent transactions....</p>}
        {res && txns.length === 0 && <p>No transactions to show</p>}
        {txns && (
          <p className="mt-2">
            You have spent a total of{' '}
            <span className="font-semibold">{total}</span> over{' '}
            <span className="font-semibold">{txns.length}</span> transactions.
          </p>
        )}
        {txns && <div className="my-4">{txns.map(transactionCard)}</div>}
      </section>
    </>
  );
}

export function formatDateForTxnCard(date: number) {
  const currentYear = Temporal.Now.zonedDateTimeISO('UTC').year;
  // the date is stored as a epoch seconds, Date constructor takes milliseconds
  return new Date(date * 1000).toLocaleDateString(['en-GB', 'ja-JP'], {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    // should be able to assume that dates without a year are from current year
    ...(epochSecToUTCYear(date) !== currentYear && {
      year: 'numeric',
    }),
  });
}

function transactionCard(txn: Transaction | SharedTxn) {
  const date = formatDateForTxnCard(txn.date);

  return (
    <article
      className="mt-3 rounded-lg border-2 border-slate-200 p-3"
      key={txn.id}
    >
      <div className="flex items-center">
        <div className="mr-4 h-10 w-10 flex-shrink-0 rounded-md bg-slate-300"></div>
        <div className="flex flex-grow">
          <div className="flex flex-grow flex-col">
            <p className="text-lg font-semibold leading-5">{txn.location}</p>
            <p className="mt-1 text-sm text-slate-500">{date}</p>
            <p className="mt-1 lowercase">
              {categoryNameFromKeyEN(txn.category)}
            </p>
            {txn.details && <p>{txn.details}</p>}
          </div>
          <p className="flex-shrink-0 text-lg font-medium text-slate-600">
            {txn.amount}
            <span className="ml-1 text-xs">円</span>
          </p>
        </div>
      </div>
    </article>
  );
}
