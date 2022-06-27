import { CategoryKey } from 'data/categories';
import { Tracker } from 'pages/shared/trackers';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import SharedTxnFormBase from './SharedTxnFormBase';

type Inputs = {
  id: string;
  shop: string;
  amount: number;
  date: string;
  settled?: boolean;
  participants: string;
  payer: string;
  category: CategoryKey;
};

async function updateSharedTxn(data: Inputs, tracker: Tracker) {
  const formData = new FormData();
  formData.append('participants', tracker.users.join(','));
  formData.append('shop', data.shop);
  formData.append('amount', data.amount.toString());
  if (!data.settled) formData.append('unsettled', 'true');
  formData.append('category', data.category);
  formData.append('payer', data.payer);

  const unixDate = new Date(data.date).getTime();
  formData.append('date', unixDate.toString());

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions/${data.id}`,
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
  txn: SharedTxn;
  tracker: Tracker;
  onApply: () => void;
  onCancel: () => void;
}

export default function SharedTxnReadUpdateForm({
  txn,
  tracker,
  onApply,
  onCancel,
}: Props) {
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, formState } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      shop: txn.shop,
      amount: txn.amount,
      date: new Date(txn.date).toISOString().split('T')[0],
      settled: !txn.unsettled,
      category: txn.category,
    },
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    data.id = txn.id;
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      updateSharedTxn(data, tracker),
    );
    onApply();
  };

  const shopInputProps = register('shop', {
    required: 'Please input a shop name',
  });
  const amountInputProps = register('amount', {
    min: {
      value: 1,
      message: 'Please input a positive amount',
    },
    required: 'Please input an amount',
  });
  const dateInputProps = register('date', {
    required: 'Please input a date',
  });
  const payerInputProps = register('payer', {
    required: 'Please select a payer',
  });
  const settledInputProps = register('settled');
  const categoryInputProps = register('category');

  return (
    <SharedTxnFormBase
      title="Update Shared Transaction"
      shopInputProps={shopInputProps}
      amountInputProps={amountInputProps}
      dateInputProps={dateInputProps}
      payerInputProps={payerInputProps}
      settledInputProps={settledInputProps}
      categoryInputProps={categoryInputProps}
      tracker={tracker}
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
    </SharedTxnFormBase>
  );
}