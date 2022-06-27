import { useUserContext } from 'context/user';
import { CategoryKey } from 'data/categories';
import { Tracker } from 'pages/shared/trackers';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import SharedTxnFormBase from './SharedTxnFormBase';

type Inputs = {
  shop: string;
  amount: number;
  date: string;
  settled?: boolean;
  participants: string;
  payer: string;
  category: CategoryKey;
};

async function createSharedTxn(data: Inputs, tracker: Tracker) {
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
  const { register, handleSubmit, setValue } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      shop: '',
      amount: 0,
      date: new Date().toISOString().split('T')[0],
      settled: false,
      payer: user.id,
      participants: '',
      category: 'unspecified.unspecified',
    },
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      createSharedTxn(data, tracker),
    );
    setValue('shop', '');
    setValue('amount', 0);
    setValue('settled', false);
    setValue('participants', '');
    setValue('category', 'unspecified.unspecified');
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
      title="Create Shared Transaction"
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
        <button className="rounded bg-indigo-500 py-2 px-4 font-bold text-white hover:bg-indigo-700 focus:outline-none focus:ring">
          Create transaction
        </button>
      </div>
    </SharedTxnFormBase>
  );
}
