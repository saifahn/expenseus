import { useUserContext } from 'context/user';
import Link from 'next/link';
import { Tracker } from 'ddb/trackers';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import { plainDateISONowString } from 'utils/dates';
import SharedTxnFormBase, {
  makeSharedTxnPayload,
  SharedTxnFormInputs,
} from './SharedTxnFormBase';

async function createSharedTxn(data: SharedTxnFormInputs, tracker: Tracker) {
  const participants = tracker.users.join(',');
  const payload = makeSharedTxnPayload({ ...data, participants });

  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
    {
      method: 'POST',
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(payload),
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
      payer: user!.id,
      participants: '',
      category: 'unspecified.unspecified',
      details: '',
      // TODO: support more than two users
      split: `${tracker.users[0]}:0.50,${tracker.users[1]}:0.50`,
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
    setValue('split', `${tracker.users[0]}:0.50,${tracker.users[1]}:0.50`);
  };

  return (
    <SharedTxnFormBase
      title="Create Shared Transaction"
      tracker={tracker}
      register={register}
      onSubmit={handleSubmit(submitCallback)}
    >
      <div className="mt-6 flex justify-end">
        <Link href={`/shared/trackers/${tracker.id}`}>
          <a className="mr-3 rounded py-2 px-4 font-medium lowercase hover:bg-slate-200 focus:outline-none focus:ring">
            Close
          </a>
        </Link>
        <button className="rounded bg-violet-500 py-2 px-4 font-medium lowercase text-white hover:bg-violet-700 focus:outline-none focus:ring">
          Create
        </button>
      </div>
    </SharedTxnFormBase>
  );
}
