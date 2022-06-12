import { useForm } from 'react-hook-form';
import { Transaction } from 'pages/personal';

interface TxnReadUpdateFormProps {
  txn: Transaction;
  onApply: () => void;
  onCancel: () => void;
}

export default function TxnReadUpdateForm({
  txn,
  onApply,
  onCancel,
}: TxnReadUpdateFormProps) {
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
