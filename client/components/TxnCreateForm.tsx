import { useUserContext } from 'context/user';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import { CategoryKey, enUSCategories } from 'data/categories';

type Inputs = {
  transactionName: string;
  amount: number;
  date: string;
  category: CategoryKey;
};

async function createTransaction(data: Inputs) {
  const formData = new FormData();
  formData.append('transactionName', data.transactionName);
  formData.append('amount', data.amount.toString());
  formData.append('category', data.category);

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

export default function TxnCreateForm() {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, setValue } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      transactionName: '',
      amount: 0,
      date: new Date().toISOString().split('T')[0],
      category: 'unspecified.unspecified',
    },
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      createTransaction(data),
    );
    setValue('transactionName', '');
    setValue('amount', 0);
    setValue('category', 'unspecified.unspecified');
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
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="text"
          id="transactionName"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold" htmlFor="amount">
          Amount
        </label>
        <input
          {...register('amount', {
            min: {
              value: 1,
              message: 'Please input a positive amount',
            },
            required: 'Please input an amount',
          })}
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
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
          className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
          type="date"
          id="date"
        />
      </div>
      <div className="mt-4">
        <label className="block font-semibold">Category</label>
        <select
          {...register('category')}
          className="mt-2 block rounded bg-white bg-clip-padding bg-no-repeat px-3 py-2 text-base font-normal text-gray-700 outline outline-1 transition ease-in-out focus:border-indigo-600 focus:bg-white focus:text-gray-700"
        >
          {enUSCategories.map((category) => (
            <option key={category.key} value={category.key}>
              {category.value}
            </option>
          ))}
        </select>
      </div>
      <div className="mt-4 flex justify-end">
        <button className="rounded bg-indigo-500 py-2 px-4 font-bold text-white hover:bg-indigo-700 focus:outline-none focus:ring">
          Create transaction
        </button>
      </div>
    </form>
  );
}
