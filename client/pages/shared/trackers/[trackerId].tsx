import { fetcher } from 'config/fetcher';
import SharedLayout from 'components/SharedLayout';
import SharedTxnSubmitForm from 'components/SharedTxnSubmitForm';
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
}

export default function TrackerPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data: tracker, error } = useSWR<Tracker>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}`,
    fetcher,
  );
  const { data: sharedTxns, error: sharedTxnsError } = useSWR<SharedTxn[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions`,
    fetcher,
  );

  return (
    <SharedLayout>
      {error && <div>Failed to load</div>}
      {!error && !tracker && <div>Loading tracker information...</div>}
      {tracker && (
        <>
          <h2 className="text-2xl mt-8">{tracker.name}</h2>
          <h3 className="mt-2">{tracker.id}</h3>
          {tracker.users.map((user) => (
            <p key={user}>{user}</p>
          ))}
        </>
      )}
      <div className="mt-8">
        <SharedTxnSubmitForm tracker={tracker} />
      </div>
      <div className="mt-8">
        <h3 className="text-2xl mt-4">Transactions</h3>
        {sharedTxnsError && <div>Failed to load</div>}
        {sharedTxns === null && <div>Loading list of transactions...</div>}
        {sharedTxns && sharedTxns.length === 0 && (
          <div>There are no transactions here yet</div>
        )}
        {sharedTxns &&
          sharedTxns.map((txn) => (
            <article className="p-2 border-2 mt-4" key={txn.id}>
              <h3 className="text-lg">{txn.shop}</h3>
              <p>{txn.amount}</p>
              <p>{txn.date}</p>
              <p>{txn.tracker}</p>
            </article>
          ))}
      </div>
    </SharedLayout>
  );
}
