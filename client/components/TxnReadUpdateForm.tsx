import { useForm, SubmitHandler } from 'react-hook-form';
import { Transaction } from 'pages/personal';
import { useSWRConfig } from 'swr';
import { useUserContext } from '../context/user';
import TxnFormBase, { createTxnFormData, TxnFormInputs } from './TxnFormBase';

async function updateTransaction(data: TxnFormInputs, txnID: string) {
  const formData = createTxnFormData(data);

  await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/${txnID}`, {
    method: 'PUT',
    headers: {
      Accept: 'application/json',
    },
    credentials: 'include',
    body: formData,
  });
}

interface Props {
  txn: Transaction;
  onApply: () => void;
  onCancel: () => void;
}

export default function TxnReadUpdateForm({ txn, onApply, onCancel }: Props) {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, formState, handleSubmit } = useForm<TxnFormInputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      location: txn.location,
      details: txn.details,
      amount: txn.amount,
      date: new Date(txn.date).toISOString().split('T')[0],
      category: txn.category,
    },
  });

  const submitCallback: SubmitHandler<TxnFormInputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      updateTransaction(data, txn.id),
    );
    onApply();
  };

  const locationInputProps = register('location', {
    required: 'Please input a location',
  });
  const amountInputProps = register('amount', {
    min: { value: 1, message: 'Please input a positive amount' },
    required: 'Please input an amount',
  });
  const dateInputProps = register('date', { required: 'Please input a date' });
  const categoryInputProps = register('category');
  const detailsInputProps = register('details');

  return (
    <TxnFormBase
      title="Update Transaction"
      locationInputProps={locationInputProps}
      amountInputProps={amountInputProps}
      dateInputProps={dateInputProps}
      categoryInputProps={categoryInputProps}
      detailsInputProps={detailsInputProps}
      onSubmit={handleSubmit(submitCallback)}
    >
      <div className="mt-4 flex justify-end">
        {formState.isDirty ? (
          <>
            <button
              className="rounded py-2 px-4 text-sm font-bold uppercase hover:bg-slate-200 focus:outline-none focus:ring"
              onClick={() => onCancel()}
            >
              Cancel
            </button>
            <button
              className="rounded bg-indigo-500 py-2 px-4 text-sm font-bold uppercase text-white hover:bg-indigo-700 focus:outline-none focus:ring"
              type="submit"
            >
              Apply
            </button>
          </>
        ) : (
          <button
            className="rounded py-2 px-4 text-sm font-bold uppercase hover:bg-slate-200 focus:outline-none focus:ring"
            onClick={() => onCancel()}
          >
            Close
          </button>
        )}
      </div>
    </TxnFormBase>
  );
}
