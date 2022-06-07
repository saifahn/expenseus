import TransactionSubmitForm from 'components/TransactionSubmitForm';
import { useUserContext } from 'context/user';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import useSWR, { mutate } from 'swr';

interface Transaction {
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

interface ReadUpdateTxnFormProps {
  txn: Transaction;
  onApply: () => void;
  onCancel: () => void;
}

function ReadUpdateTxnForm({ txn, onApply, onCancel }: ReadUpdateTxnFormProps) {
  const { register, formState } = useForm({
    shouldUseNativeValidation: true,
    defaultValues: {
      transactionName: txn.name,
      amount: txn.amount,
      date: new Date(txn.date).toISOString().split('T')[0],
      image: null,
    },
  });

  return (
    <form onSubmit={(e) => e.preventDefault()} className="border-4 p-6 mt-4">
      <h3 className="text-lg font-semibold">Update Transaction</h3>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="name">
          Name
        </label>
        <input
          {...register('transactionName', {
            required: 'Please input a transaction name',
          })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="text"
          id="transactionName"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...register('amount', { required: 'Please input an amount' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="text"
          inputMode="numeric"
          pattern="[0-9]*"
          id="amount"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="date">
          Date
        </label>
        <input
          {...register('date', { required: 'Please input a date' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="date"
          id="date"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="addPicture">
          Add a picture?
        </label>
        <input
          {...register('image')}
          id="addPicture"
          type="file"
          role="button"
          aria-label="Add picture"
          accept="image/*"
        />
      </div>
      <div className="mt-4 flex justify-end">
        {formState.isDirty ? (
          <>
            <button
              className="hover:bg-slate-200 font-bold uppercase text-sm py-2 px-4 rounded focus:outline-none focus:ring"
              onClick={() => onCancel()}
            >
              Cancel
            </button>
            <button
              className="bg-indigo-500 hover:bg-indigo-700 text-white font-bold uppercase text-sm py-2 px-4 rounded focus:outline-none focus:ring"
              onClick={() => onApply()}
            >
              Apply
            </button>
          </>
        ) : (
          <button
            className="hover:bg-slate-200 font-bold uppercase text-sm py-2 px-4 rounded focus:outline-none focus:ring"
            onClick={() => onCancel()}
          >
            Close
          </button>
        )}
      </div>
    </form>
  );
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
        <ReadUpdateTxnForm
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
