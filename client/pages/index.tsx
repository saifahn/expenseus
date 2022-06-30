import { useUserContext } from 'context/user';
import useSWR from 'swr';
import { Transaction } from 'types/Transaction';
import { Temporal } from 'temporal-polyfill';
import {
  epochSecToLocaleString,
  plainDateStringToEpochSec,
} from 'utils/temporal';

export default function Home() {
  const { user } = useUserContext();
  const todayEpochSeconds = plainDateStringToEpochSec(
    Temporal.Now.plainDateISO().toString(),
  );
  const threeMonthsAgoEpochSeconds = plainDateStringToEpochSec(
    Temporal.Now.plainDateISO().subtract({ months: 3 }).toString(),
  );
  const { data: txns, error } = useSWR<Transaction[]>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}/range?from=${threeMonthsAgoEpochSeconds}&to=${todayEpochSeconds}`
      : null,
  );

  return (
    <>
      <div>Hi, {user.username}!</div>
      {error && <div>Failed to load recent transactions</div>}
      {txns === null && <div>Loading recent transactions....</div>}
      {txns && txns.length === 0 && <div>No transactions to show</div>}
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
