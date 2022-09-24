import TrackerLayout from 'components/LayoutTracker';
import SharedTxnCreateForm from 'components/SharedTxnCreateForm';
import SharedTxnReadUpdateForm from 'components/SharedTxnReadUpdateForm';
import { categoryNameFromKeyEN, SubcategoryKey } from 'data/categories';
import Link from 'next/link';
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

/**
 * SharedTxnOne displays one shared transaction to be showed in a list.
 */
type SharedTxnOneProps = {
  txn: SharedTxn;
  onTxnClick: (txn: SharedTxn) => void;
};

function SharedTxnOne({ txn, onTxnClick }: SharedTxnOneProps) {
  return (
    <article
      className="mt-4 cursor-pointer border-2 p-2 hover:bg-slate-200 active:bg-slate-300"
      onClick={() => onTxnClick(txn)}
      key={txn.id}
    >
      <div className="flex justify-between">
        <h3 className="text-lg">{txn.location}</h3>
      </div>
      <p>
        {txn.amount} paid by {txn.payer}
      </p>
      <p>{categoryNameFromKeyEN(txn.category)}</p>
      <p>{epochSecToLocaleString(txn.date)}</p>
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
      <Link href={`/shared/trackers/${trackerId}/create-txn`}>
        <a className="mt-4 block rounded-lg bg-violet-50 p-3 font-medium lowercase text-black hover:bg-violet-100 active:bg-violet-200">
          Create new transaction +
        </a>
      </Link>
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
            <div className="mt-6">
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
                  />
                ))}
            </div>
          </>
        )
      )}
    </TrackerLayout>
  );
}
