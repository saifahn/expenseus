import { useForm, SubmitHandler } from 'react-hook-form';
import { Transaction } from 'pages/personal';
import { useSWRConfig } from 'swr';
import { useUserContext } from '../context/user';

interface TxnReadUpdateFormProps {
  txn: Transaction;
  onApply: () => void;
  onCancel: () => void;
}

type Inputs = {
  txnID: string;
  transactionName: string;
  amount: number;
  date: string;
};

async function updateTransaction(data: Inputs) {
  const formData = new FormData();
  formData.append('transactionName', data.transactionName);
  formData.append('amount', data.amount.toString());

  const unixDate = new Date(data.date).getTime();
  formData.append('date', unixDate.toString());

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/${data.txnID}`,
    {
      method: 'PUT',
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include',
      body: formData,
    },
  );
}

export default function TxnReadUpdateForm({
  txn,
  onApply,
  onCancel,
}: TxnReadUpdateFormProps) {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, formState, handleSubmit } = useForm({
    shouldUseNativeValidation: true,
    defaultValues: {
      transactionName: txn.name,
      amount: txn.amount,
      date: new Date(txn.date).toISOString().split('T')[0],
    },
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    data.txnID = txn.id;
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      updateTransaction(data),
    );
    onApply();
  };

  return (
    <form onSubmit={handleSubmit(submitCallback)} className="border-4 p-6 mt-4">
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
              type="submit"
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
