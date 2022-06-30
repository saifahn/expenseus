import TxnCreateForm from 'components/TxnCreateForm';
import TxnReadUpdateForm from 'components/TxnReadUpdateForm';
import { useUserContext } from 'context/user';
import { useState } from 'react';
import useSWR, { mutate } from 'swr';
import { epochSecToLocaleString } from 'utils/temporal';
import { Transaction } from 'types/Transaction';
import PersonalLayout from 'components/LayoutPersonal';

async function deleteTransaction(txnId: string) {
  await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/${txnId}`, {
    method: 'DELETE',
    headers: {
      Accept: 'application/json',
    },
    credentials: 'include',
  });
}

type TxnOneProps = {
  txn: Transaction;
  onTxnClick: (txn: Transaction) => void;
};

function TxnOne({ txn, onTxnClick }: TxnOneProps) {
  const { user } = useUserContext();

  function handleDelete(e: React.MouseEvent) {
    e.stopPropagation();
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      deleteTransaction(txn.id),
    );
  }

  return (
    <article
      className="mt-4 cursor-pointer border-2 p-2 hover:bg-slate-200 active:bg-slate-300"
      key={txn.id}
      onClick={() => onTxnClick(txn)}
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
      <p>{txn.amount}</p>
      <p>{txn.category}</p>
      {txn.details && <p>{txn.details}</p>}
      <p>{epochSecToLocaleString(txn.date)}</p>
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
        {selectedTxn ? (
          <TxnReadUpdateForm
            txn={selectedTxn}
            onApply={() => setSelectedTxn(null)}
            onCancel={() => setSelectedTxn(null)}
          />
        ) : (
          <>
            <TxnCreateForm />
            <div className="mt-4 p-4">
              <h2 className="text-2xl">Personal transactions</h2>
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
          </>
        )}
      </PersonalLayout>
    </>
  );
}
