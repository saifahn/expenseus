import { Tracker } from 'pages/shared/trackers';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import SharedTxnFormBase, {
  createSharedTxnFormData,
  SharedTxnFormInputs,
} from './SharedTxnFormBase';

async function updateSharedTxn(
  data: SharedTxnFormInputs,
  tracker: Tracker,
  txnID: string,
) {
  const formData = createSharedTxnFormData(data);
  formData.append('participants', tracker.users.join(','));

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions/${txnID}`,
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
  const { register, handleSubmit, formState } = useForm<SharedTxnFormInputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      location: txn.location,
      amount: txn.amount,
      date: new Date(txn.date).toISOString().split('T')[0],
      settled: !txn.unsettled,
      category: txn.category,
      details: txn.details,
    },
  });

  const submitCallback: SubmitHandler<SharedTxnFormInputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      updateSharedTxn(data, tracker, txn.id),
    );
    onApply();
  };

  return (
    <SharedTxnFormBase
      title="Update Shared Transaction"
      tracker={tracker}
      register={register}
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
