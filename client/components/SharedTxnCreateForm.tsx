import { useUserContext } from 'context/user';
import { Tracker } from 'pages/shared/trackers';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import { plainDateISONowString } from 'utils/temporal';
import SharedTxnFormBase, {
  createSharedTxnFormData,
  SharedTxnFormInputs,
} from './SharedTxnFormBase';

async function createSharedTxn(data: SharedTxnFormInputs, tracker: Tracker) {
  const formData = createSharedTxnFormData(data);
  formData.append('participants', tracker.users.join(','));

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
    {
      method: 'POST',
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include',
      body: formData,
    },
  );
}

interface Props {
  tracker: Tracker;
}

export default function SharedTxnCreateForm({ tracker }: Props) {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, setValue } = useForm<SharedTxnFormInputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      location: '',
      amount: 0,
      date: plainDateISONowString(),
      settled: false,
      payer: user.id,
      participants: '',
      category: 'unspecified.unspecified',
      details: '',
    },
  });

  const submitCallback: SubmitHandler<SharedTxnFormInputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      createSharedTxn(data, tracker),
    );
    setValue('location', '');
    setValue('amount', 0);
    setValue('settled', false);
    setValue('participants', '');
    setValue('category', 'unspecified.unspecified');
    setValue('details', '');
  };

  return (
    <SharedTxnFormBase
      title="Create Shared Transaction"
      tracker={tracker}
      register={register}
      onSubmit={handleSubmit(submitCallback)}
    >
      <div className="mt-4 flex justify-end">
        <button className="rounded bg-indigo-500 py-2 px-4 font-bold text-white hover:bg-indigo-700 focus:outline-none focus:ring">
          Create transaction
        </button>
      </div>
    </SharedTxnFormBase>
  );
}
