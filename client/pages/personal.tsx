import TransactionSubmitForm from 'components/TransactionSubmitForm';
import TxnReadUpdateForm from 'components/TxnReadUpdateForm';
import { useUserContext } from 'context/user';
import { useState } from 'react';
import useSWR, { mutate } from 'swr';

export interface Transaction {
  name: string;
  id: string;
  userId: string;
  amount: number;
  imageUrl?: string;
  date: string;
}

async function deleteTransaction(txnId: string) {
  await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/${txnId}`, {
    method: 'DELETE',
    headers: {
      Accept: 'application/json',
    },
    credentials: 'include',
  });
}

export default function Personal() {
  const { user } = useUserContext();
  const [selectedTxn, setSelectedTxn] = useState<Transaction | null>(null);
  const { data: transactions, error } = useSWR<Transaction[]>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`
      : null,
  );

  function handleDelete(txnId: string) {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      deleteTransaction(txnId),
    );
  }

  return (
    <>
      <h1 className="text-4xl">Personal</h1>
      {selectedTxn ? (
        <TxnReadUpdateForm
          txn={selectedTxn}
          onApply={() => setSelectedTxn(null)}
          onCancel={() => setSelectedTxn(null)}
        />
      ) : (
        <>
          <div className="mt-4">
            <TransactionSubmitForm />
          </div>
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
                <article
                  className="p-2 border-2 mt-4 hover:bg-slate-200 active:bg-slate-300 cursor-pointer"
                  key={txn.id}
                  onClick={() => setSelectedTxn(txn)}
                >
                  <div className="flex justify-between">
                    <h3 className="text-lg">{txn.name}</h3>
                    <button
                      className="bg-red-500 hover:bg-red-700 active:bg-blue-300 text-white font-bold uppercase text-sm py-2 px-4 rounded focus:outline-none focus:ring"
                      onClick={() => handleDelete(txn.id)}
                    >
                      Delete
                    </button>
                  </div>
                  <p>{txn.amount}</p>
                  <p>{txn.userId}</p>
                  <p>{txn.id}</p>
                  <p>{new Date(txn.date).toDateString()}</p>
                </article>
              ))}
          </div>
        </>
      )}
    </>
  );
}
