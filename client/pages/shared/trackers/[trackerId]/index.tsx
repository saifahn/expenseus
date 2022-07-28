import TrackerLayout from 'components/LayoutTracker';
import SharedTxnCreateForm from 'components/SharedTxnCreateForm';
import SharedTxnReadUpdateForm from 'components/SharedTxnReadUpdateForm';
import { SubcategoryKey } from 'data/categories';
import { useRouter } from 'next/router';
import React, { useState } from 'react';
import useSWR, { mutate } from 'swr';
import { epochSecToLocaleString } from 'utils/dates';
import { Tracker } from '..';

export interface SharedTxn {
  id: string;
  location: string;
  amount: number;
  date: number;
  participants: string[];
  tracker: string;
  category: SubcategoryKey;
  payer: string;
  details: string;
  unsettled?: boolean;
  split?: {
    [userId: string]: number;
  };
}

async function deleteSharedTxn(txn: SharedTxn) {
  const payload = {
    tracker: txn.tracker,
    txnID: txn.id,
    participants: txn.participants,
  };

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${txn.tracker}/transactions/${txn.id}`,
    {
      method: 'DELETE',
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(payload),
    },
  );
}

/**
 * SharedTxnOne displays one shared transaction to be showed in a list.
 */
type SharedTxnOneProps = {
  txn: SharedTxn;
  tracker: Tracker;
  onTxnClick: (txn: SharedTxn) => void;
};

function SharedTxnOne({ txn, tracker, onTxnClick }: SharedTxnOneProps) {
  function handleDelete(e: React.MouseEvent) {
    e.stopPropagation();
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      deleteSharedTxn(txn),
    );
  }

  return (
    <article
      className="mt-4 cursor-pointer border-2 p-2 hover:bg-slate-200 active:bg-slate-300"
      onClick={() => onTxnClick(txn)}
      key={txn.id}
    >
      <div className="flex justify-between">
        <h3 className="text-lg">{txn.location}</h3>
        <button
          className="rounded bg-red-500 py-2 px-4 text-sm font-bold uppercase text-white hover:bg-red-700 focus:outline-none focus:ring active:bg-blue-300"
          onClick={handleDelete}
        >
          Delete
        </button>
      </div>
      <p>
        {txn.amount} paid by {txn.payer}
      </p>
      <p>{txn.category}</p>
      <p>{epochSecToLocaleString(txn.date)}</p>
      <p>{txn.tracker}</p>
      {txn.details && <p>{txn.details}</p>}
    </article>
  );
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
    <TrackerLayout>
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
        tracker && (
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
                  <SharedTxnOne
                    txn={txn}
                    onTxnClick={setSelectedTxn}
                    key={txn.id}
                    tracker={tracker}
                  />
                ))}
            </div>
          </>
        )
      )}
    </TrackerLayout>
  );
}
