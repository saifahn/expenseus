import { useUserContext } from 'context/user';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';

type Inputs = {
  transactionName: string;
  amount: number;
  date: string;
  image: File;
};

async function createTransaction(data: Inputs) {
  const formData = new FormData();
  formData.append('transactionName', data.transactionName);
  formData.append('amount', data.amount.toString());
  formData.append('image', data.image);

  const unixDate = new Date(data.date).getTime();
  formData.append('date', unixDate.toString());

  await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions`, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
    },
    credentials: 'include',
    body: formData,
  });
}

export default function TransactionSubmitForm() {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, setValue } = useForm({
    shouldUseNativeValidation: true,
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      createTransaction(data),
    );
    setValue('transactionName', '');
    setValue('amount', 0);
    setValue('image', null);
  };

  return (
    <form onSubmit={handleSubmit(submitCallback)} className="border-4 p-6">
      <h3 className="text-lg font-semibold">Create Transaction</h3>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="name">
          Name
        </label>
        <input
          {...register('transactionName', {
            required: 'Please input a transaction name',
          })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="text"
          id="transactionName"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...register('amount', { required: 'Please input an amount' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="text"
          inputMode="numeric"
          pattern="[0-9]*"
          id="amount"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="date">
          Date
        </label>
        <input
          {...register('date', { required: 'Please input a date' })}
          className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
          type="date"
          id="date"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="addPicture">
          Add a picture?
        </label>
        <input
          {...register('image')}
          id="addPicture"
          type="file"
          role="button"
          aria-label="Add picture"
          accept="image/*"
        />
      </div>
      <div className="mt-4 flex justify-end">
        <button className="bg-indigo-500 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:ring">
          Create transaction
        </button>
      </div>
    </form>
  );
}
