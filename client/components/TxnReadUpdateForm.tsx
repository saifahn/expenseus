import { useForm, SubmitHandler } from 'react-hook-form';
import { Transaction } from 'types/Transaction';
import { useSWRConfig } from 'swr';
import { useUserContext } from '../context/user';
import TxnFormBase, {
  makeCreateTxnPayload,
  TxnFormInputs,
} from './TxnFormBase';
import { epochSecToISOString } from 'utils/dates';

async function updateTransaction(data: TxnFormInputs, txnID: string) {
  await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/${txnID}`, {
    method: 'PUT',
    headers: {
      Accept: 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(makeCreateTxnPayload(data)),
  });
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
      date: epochSecToISOString(txn.date),
      category: txn.category,
    },
  });

  const submitCallback: SubmitHandler<TxnFormInputs> = (data) => {
    if (!user) {
      console.error('user not loaded');
      return;
    }
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      updateTransaction(data, txn.id),
    );
    onApply();
  };

  function handleDelete(e: React.MouseEvent) {
    if (!user) {
      console.error('user not loaded');
      return;
    }
    e.stopPropagation();
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      deleteTransaction(txn.id),
    );
  }

  return (
    <TxnFormBase
      title="Update Transaction"
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
              className="mr-3 rounded py-2 px-4 font-medium lowercase hover:bg-slate-200 focus:outline-none focus:ring"
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
    </TxnFormBase>
  );
}
