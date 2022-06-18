import SharedLayout from 'components/SharedLayout';
import SharedTxnCreateForm from 'components/SharedTxnCreateForm';
import { useRouter } from 'next/router';
import useSWR from 'swr';
import { Tracker } from '.';

interface SharedTxn {
  id: string;
  shop: string;
  amount: number;
  date: string;
  unsettled?: boolean;
  tracker: string;
  category: string;
}

export default function TrackerPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data: tracker, error } = useSWR<Tracker>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}`,
  );
  const { data: sharedTxns, error: sharedTxnsError } = useSWR<SharedTxn[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions`,
  );

  return (
    <SharedLayout>
      {error && <div>Failed to load</div>}
      {!error && !tracker && <div>Loading tracker information...</div>}
      {tracker && (
        <>
          <h2 className="mt-8 text-2xl">{tracker.name}</h2>
          <h3 className="mt-2">{tracker.id}</h3>
          {tracker.users.map((user) => (
            <p key={user}>{user}</p>
          ))}
          <div className="mt-8">
            <SharedTxnCreateForm tracker={tracker} />
          </div>
        </>
      )}

      <div className="mt-8">
        <h3 className="mt-4 text-2xl">Transactions</h3>
        {sharedTxnsError && <div>Failed to load</div>}
        {sharedTxns === null && <div>Loading list of transactions...</div>}
        {sharedTxns && sharedTxns.length === 0 && (
          <div>There are no transactions here yet</div>
        )}
        {sharedTxns &&
          sharedTxns.map((txn) => (
            <article className="mt-4 border-2 p-2" key={txn.id}>
              <h3 className="text-lg">{txn.shop}</h3>
              <p>{txn.amount}</p>
              <p>{txn.category}</p>
              <p>{new Date(txn.date).toDateString()}</p>
              <p>{txn.tracker}</p>
            </article>
          ))}
      </div>
    </SharedLayout>
  );
}
