import { Tracker } from 'pages/shared/trackers';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';

type Inputs = {
  // file: File[];
  shop: string;
  amount: number;
  date: string;
  unsettled?: boolean;
  participants: string;
};

async function createSharedTxn(data: Inputs, tracker: Tracker) {
  const formData = new FormData();
  formData.append('participants', tracker.users.join(','));
  formData.append('shop', data.shop);
  formData.append('amount', data.amount.toString());
  if (data.unsettled) formData.append('unsettled', 'true');

  const unixDate = new Date(data.date).getTime();
  formData.append('date', unixDate.toString());
  // make formData
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

export default function SharedTxnSubmitForm({ tracker }: Props) {
  const { mutate } = useSWRConfig();
  const { register, handleSubmit } = useForm({
    shouldUseNativeValidation: true,
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${tracker.id}/transactions`,
      createSharedTxn(data, tracker),
    );
  };

  return (
    <form onSubmit={handleSubmit(submitCallback)} className="border-4 p-6">
      <h3 className="text-lg font-semibold">Create Transaction</h3>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="shop">
          Shop
        </label>
        <input
          {...register('shop', { required: 'Please input the shop name' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="text"
          id="shop"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...register('amount', { required: 'Please input the amount' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="number"
          id="amount"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="date">
          Date
        </label>
        <input
          {...register('date', { required: 'Please input the date' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="date"
          id="date"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="unsettled">
          Unsettled?
        </label>
        <input {...register('unsettled')} type="checkbox" id="unsettled" />
      </div>
      <div className="mt-4 flex justify-end">
        <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:ring">
          Create transaction
        </button>
      </div>
    </form>
  );
}
