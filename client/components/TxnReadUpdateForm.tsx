import { useForm, SubmitHandler } from 'react-hook-form';
import { Transaction } from 'pages/personal';
import { useSWRConfig } from 'swr';
import { useUserContext } from '../context/user';
import { CategoryKey } from 'data/categories';
import TxnFormBase from './TxnFormBase';

type Inputs = {
  txnID: string;
  location: string;
  amount: number;
  date: string;
  category: CategoryKey;
  details: string;
};

async function updateTransaction(data: Inputs) {
  const formData = new FormData();
  formData.append('location', data.location);
  formData.append('details', data.details);
  formData.append('amount', data.amount.toString());
  formData.append('category', data.category);

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

interface Props {
  txn: Transaction;
  onApply: () => void;
  onCancel: () => void;
}

export default function TxnReadUpdateForm({ txn, onApply, onCancel }: Props) {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, formState, handleSubmit } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      location: txn.location,
      details: txn.details,
      amount: txn.amount,
      date: new Date(txn.date).toISOString().split('T')[0],
      category: txn.category,
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
