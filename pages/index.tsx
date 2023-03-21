import { useUserContext } from 'context/user';
import useSWR from 'swr';
import { Transaction } from 'types/Transaction';
import { Temporal } from 'temporal-polyfill';
import { plainDateStringToEpochSec } from 'utils/dates';
import { SharedTxn } from './shared/trackers/[trackerId]';
import { calculatePersonalTotal } from 'utils/analysis';
import { categoryNameFromKeyEN, getEmojiForTxnCard } from 'data/categories';
import { formatDateForTxnCard } from 'utils/dates';
import Link from 'next/link';
import Head from 'next/head';

export type AllTxnsResponse = {
  transactions: Transaction[];
  sharedTransactions: SharedTxn[];
};

function returnSortedTxnsAndPersonalTotal(
  allTxnsRes?: AllTxnsResponse,
  id?: string,
) {
  let txns: (Transaction | SharedTxn)[] = [];
  if (allTxnsRes) {
    txns = [...allTxnsRes.transactions, ...allTxnsRes.sharedTransactions].sort(
      (a, b) => b.date - a.date, // sorted descending, so the most recent dates are shown first
    );
  }
  let total;
  if (id) {
    total = calculatePersonalTotal(id, txns);
  }
  return { txns, total };
}

export default function Home() {
  const { user, error: userError } = useUserContext();
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

  const { txns, total } = returnSortedTxnsAndPersonalTotal(res, user?.id);

  return (
    <>
      <Head>
        <title>home - expenseus</title>
      </Head>
      <nav className="mt-4">
        <ul className="flex">
          <li className="mr-4 flex flex-1">
            <Link href="/all-personal-analysis">
              <a className="w-full rounded-lg bg-emerald-50 py-3 px-4 font-medium lowercase text-black hover:bg-emerald-100 active:bg-emerald-200">
                🔎 Analyze
              </a>
            </Link>
          </li>
        </ul>
      </nav>
      <section>
        {userError && <p>Failed to load user</p>}
        {user && <p className="mt-4">Hi, {user.username}!</p>}
        {error && <p>Failed to load recent transactions</p>}
        {res === null && <p>Loading recent transactions....</p>}
        {res && txns.length === 0 && <p>No transactions to show.</p>}
        {txns.length > 0 && (
          <p className="mt-2">
            In the last 90 days, you have spent a total of{' '}
            <span className="font-semibold">{total}</span> over{' '}
            <span className="font-semibold">{txns.length}</span> transactions.
          </p>
        )}
        {txns && <div className="my-4">{txns.map(transactionCard)}</div>}
      </section>
    </>
  );
}

function transactionCard(txn: Transaction | SharedTxn) {
  const emoji = getEmojiForTxnCard(txn.category);
  const date = formatDateForTxnCard(txn.date);

  return (
    <article
      className="mt-3 rounded-lg border-2 border-slate-200 p-3"
      key={txn.id}
    >
      <div className="flex items-center">
        <div className="mr-3 flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-xl">
          {emoji}
        </div>
        <div className="flex flex-grow">
          <div className="flex flex-grow flex-col pr-2">
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
