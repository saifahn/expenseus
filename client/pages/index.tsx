import { useUserContext } from 'context/user';
import useSWR from 'swr';
import { Transaction } from 'types/Transaction';
import { Temporal } from 'temporal-polyfill';
import { epochSecToLocaleString, plainDateStringToEpochSec } from 'utils/dates';
import { SharedTxn } from './shared/trackers/[trackerId]';
import { calculatePersonalTotal } from 'utils/analysis';

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
      <div>Hi, {user.username}!</div>
      {error && <div>Failed to load recent transactions</div>}
      {res === null && <div>Loading recent transactions....</div>}
      {res && txns.length === 0 && <div>No transactions to show</div>}
      {txns && (
        <p className="mt-4">
          You have spent a total of {total} over {txns.length} transactions.
        </p>
      )}
      {txns?.map((txn) => (
        <article
          className="mt-4 cursor-pointer border-2 p-2 hover:bg-slate-200 active:bg-slate-300"
          key={txn.id}
        >
          <div className="flex justify-between">
            <h3 className="text-lg">{txn.location}</h3>
          </div>
          <p>{txn.amount}</p>
          <p>{txn.category}</p>
          {txn.details && <p>{txn.details}</p>}
          <p>{epochSecToLocaleString(txn.date)}</p>
        </article>
      ))}
    </>
  );
}
