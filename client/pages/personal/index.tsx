import TxnReadUpdateForm from 'components/TxnReadUpdateForm';
import { useUserContext } from 'context/user';
import { useState } from 'react';
import useSWR from 'swr';
import { Transaction } from 'types/Transaction';
import PersonalLayout from 'components/LayoutPersonal';
import { categoryNameFromKeyEN, getEmojiForTxnCard } from 'data/categories';
import { formatDateForTxnCard } from 'pages';

type TxnOneProps = {
  txn: Transaction;
  onTxnClick: (txn: Transaction) => void;
};

function TxnOne({ txn, onTxnClick }: TxnOneProps) {
  const emoji = getEmojiForTxnCard(txn.category);
  const date = formatDateForTxnCard(txn.date);

  return (
    <article
      className="mt-3 cursor-pointer rounded-lg border-2 border-slate-200 p-3 hover:bg-slate-200 active:bg-slate-300"
      key={txn.id}
      onClick={() => onTxnClick(txn)}
    >
      <div className="flex items-center">
        <div className="mr-4 flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-xl">
          {emoji}
        </div>
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

export default function Personal() {
  const { user } = useUserContext();
  const [selectedTxn, setSelectedTxn] = useState<Transaction>(null);
  const { data: transactions, error } = useSWR<Transaction[]>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`
      : null,
  );

  return (
    <>
      <PersonalLayout>
        <div className="relative pb-5">
          {!selectedTxn && (
            <div>
              {error && <div>Failed to load transactions</div>}
              {transactions === null && (
                <div>Loading list of transactions...</div>
              )}
              {transactions && transactions.length === 0 && (
                <div>No transactions to show</div>
              )}
              {transactions &&
                transactions.map((txn) => (
                  <TxnOne txn={txn} onTxnClick={setSelectedTxn} key={txn.id} />
                ))}
            </div>
          )}
          <div
            className={[
              'absolute top-0 w-full transition-all',
              selectedTxn ? 'opacity-100' : 'opacity-0',
            ].join(' ')}
          >
            {selectedTxn && (
              <TxnReadUpdateForm
                txn={selectedTxn}
                onApply={() => setSelectedTxn(null)}
                onCancel={() => setSelectedTxn(null)}
              />
            )}
          </div>
        </div>
      </PersonalLayout>
    </>
  );
}
