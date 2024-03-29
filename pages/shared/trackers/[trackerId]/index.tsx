import TrackerLayout from 'components/LayoutTracker';
import SharedTxnReadUpdateForm from 'components/SharedTxnReadUpdateForm';
import {
  categoryNameFromKeyEN,
  getEmojiForTxnCard,
  SubcategoryKey,
} from 'data/categories';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { formatDateForTxnCard } from 'utils/dates';
import { useState } from 'react';
import useSWR from 'swr';
import { Tracker } from 'ddb/trackers';
import Head from 'next/head';

export type SharedTxn = {
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
};

/**
 * SharedTxnOne displays one shared transaction to be showed in a list.
 */
type SharedTxnOneProps = {
  txn: SharedTxn;
  onTxnClick?: (txn: SharedTxn) => void;
};

export function SharedTxnOne({ txn, onTxnClick }: SharedTxnOneProps) {
  const emoji = getEmojiForTxnCard(txn.category);
  const date = formatDateForTxnCard(txn.date);
  function handleClick(txn: SharedTxn) {
    if (onTxnClick) {
      onTxnClick(txn);
    }
    return;
  }

  return (
    <article
      className={`mt-3  rounded-lg border-2 border-slate-200 p-3 ${
        onTxnClick && 'cursor-pointer hover:bg-slate-200 active:bg-slate-300'
      }`}
      key={txn.id}
      onClick={() => handleClick(txn)}
    >
      <div className="flex items-center">
        <div className="mx-2 h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-xl">
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
            <p className="mt-1 text-sm text-slate-400">{txn.payer}</p>
          </div>
          <div className="flex flex-shrink-0 flex-col items-end">
            <p className="text-lg font-medium text-slate-600">
              {txn.amount}
              <span className="ml-1 text-xs">円</span>
            </p>
          </div>
        </div>
      </div>
    </article>
  );
}

export default function TrackerPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const [selectedTxn, setSelectedTxn] = useState<SharedTxn | null>(null);
  const { data: tracker, error } = useSWR<Tracker>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}`,
  );
  const { data: sharedTxns, error: sharedTxnsError } = useSWR<SharedTxn[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions`,
  );

  return (
    <TrackerLayout>
      <Head>
        {
          <title>
            {`${selectedTxn ? 'update shared transaction' : tracker?.name} - `}
            expenseus
          </title>
        }
      </Head>
      <div className="relative pb-5">
        {!selectedTxn && (
          <>
            <Link href={`/shared/trackers/${trackerId}/create-txn`}>
              <a className="mt-4 block rounded-lg bg-violet-50 p-3 font-medium lowercase text-black hover:bg-violet-100 active:bg-violet-200">
                ➕ Create new transaction
              </a>
            </Link>
            <div className="mt-2 flex">
              <Link href={`/shared/trackers/${trackerId}/unsettled-txns`}>
                <a className="mr-4 w-1/2 rounded-lg bg-rose-50 py-3 px-4 font-medium lowercase text-black hover:bg-rose-100  active:bg-rose-200">
                  See unsettled
                </a>
              </Link>
              <Link href={`/shared/trackers/${trackerId}/analysis`}>
                <a className="w-1/2 rounded-lg bg-emerald-50 py-3 px-4 font-medium lowercase text-black hover:bg-emerald-100 active:bg-emerald-200">
                  🔎 Analyze
                </a>
              </Link>
            </div>
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
        )}
        <div
          className={[
            'absolute top-0 w-full transition-all',
            selectedTxn ? 'opacity-100' : 'opacity-0',
          ].join(' ')}
        >
          {selectedTxn && tracker && (
            <SharedTxnReadUpdateForm
              txn={selectedTxn}
              tracker={tracker}
              onApply={() => setSelectedTxn(null)}
              onCancel={() => setSelectedTxn(null)}
            />
          )}
        </div>
      </div>
    </TrackerLayout>
  );
}
