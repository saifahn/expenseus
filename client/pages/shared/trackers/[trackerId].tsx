import SharedLayout from 'components/SharedLayout';
import SharedTxnCreateForm from 'components/SharedTxnCreateForm';
import SharedTxnReadUpdateForm from 'components/SharedTxnReadUpdateForm';
import { CategoryKey } from 'data/categories';
import { useRouter } from 'next/router';
import { useState } from 'react';
import useSWR from 'swr';
import { Tracker } from '.';

export interface SharedTxn {
  id: string;
  shop: string;
  amount: number;
  date: string;
  unsettled?: boolean;
  tracker: string;
  category: CategoryKey;
}

export default function TrackerPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const [selectedTxn, setSelectedTxn] = useState<SharedTxn>(null);
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
          {selectedTxn ? (
            <div className="mt-4">
              <SharedTxnReadUpdateForm
                txn={selectedTxn}
                tracker={tracker}
                onApply={() => setSelectedTxn(null)}
                onCancel={() => setSelectedTxn(null)}
              />
            </div>
          ) : (
            <>
              <div className="mt-4">
                <SharedTxnCreateForm tracker={tracker} />
              </div>
              <div className="mt-8">
                <h3 className="mt-4 text-2xl">Transactions</h3>
                {sharedTxnsError && <div>Failed to load</div>}
                {sharedTxns === null && (
                  <div>Loading list of transactions...</div>
                )}
                {sharedTxns && sharedTxns.length === 0 && (
                  <div>There are no transactions here yet</div>
                )}
                {sharedTxns &&
                  sharedTxns.map((txn) => (
                    <article
                      className="mt-4 cursor-pointer border-2 p-2 hover:bg-slate-200 active:bg-slate-300"
                      onClick={() => setSelectedTxn(txn)}
                      key={txn.id}
                    >
                      <h3 className="text-lg">{txn.shop}</h3>
                      <p>{txn.amount}</p>
                      <p>{txn.category}</p>
                      <p>{new Date(txn.date).toDateString()}</p>
                      <p>{txn.tracker}</p>
                    </article>
                  ))}
              </div>
            </>
          )}
        </>
      )}
    </SharedLayout>
  );
}
