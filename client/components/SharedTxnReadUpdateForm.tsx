import { Tracker } from 'pages/shared/trackers';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import { epochSecToISOString } from 'utils/dates';
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

async function deleteSharedTxn(txn: SharedTxn) {
  const payload = {
    tracker: txn.tracker,
    txnID: txn.id,
    participants: txn.participants,
  };

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${txn.tracker}/transactions/${txn.id}`,
    {
      method: 'DELETE',
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(payload),
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

  let userSplits = [];
  for (const userSplit of Object.entries(txn.split)) {
    userSplits.push(`${userSplit[0]}:${userSplit[1]}`);
  }

  const { register, handleSubmit, formState } = useForm<SharedTxnFormInputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      location: txn.location,
      amount: txn.amount,
      date: epochSecToISOString(txn.date),
      settled: !txn.unsettled,
      category: txn.category,
      details: txn.details,
      payer: txn.payer,
      split: userSplits.join(','),
    },
  });

  const submitCallback: SubmitHandler<SharedTxnFormInputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      updateSharedTxn(data, tracker, txn.id),
    );
    onApply();
  };

  function handleDelete(e: React.MouseEvent) {
    e.stopPropagation();
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      deleteSharedTxn(txn),
    );
  }

  return (
    <SharedTxnFormBase
      title="Update Shared Transaction"
      tracker={tracker}
      register={register}
      onSubmit={handleSubmit(submitCallback)}
    >
      <div className="mt-6 flex">
        <div className="flex-grow">
          <button
            className="rounded bg-red-500 py-2 px-4 font-medium lowercase text-white hover:bg-red-700 focus:outline-none focus:ring active:bg-red-300"
            onClick={handleDelete}
          >
            Delete
          </button>
        </div>
        {formState.isDirty ? (
          <>
            <button
              className="mr-2 rounded py-2 px-4 font-medium lowercase hover:bg-slate-200 focus:outline-none focus:ring"
              onClick={onCancel}
            >
              Cancel
            </button>
            <button
              className="rounded bg-violet-500 py-2 px-4 font-medium lowercase text-white hover:bg-violet-700 focus:outline-none focus:ring"
              type="submit"
            >
              Apply
            </button>
          </>
        ) : (
          <button
            className="rounded py-2 px-4 font-medium lowercase hover:bg-slate-200 focus:outline-none focus:ring"
            onClick={onCancel}
          >
            Close
          </button>
        )}
      </div>
    </SharedTxnFormBase>
  );
}
