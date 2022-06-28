import { useUserContext } from 'context/user';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import { CategoryKey } from 'data/categories';
import TxnFormBase from './TxnFormBase';

type Inputs = {
  location: string;
  amount: number;
  date: string;
  category: CategoryKey;
  details: string;
};

async function createTransaction(data: Inputs) {
  const formData = new FormData();
  formData.append('location', data.location);
  formData.append('details', data.details);
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
      location: '',
      details: '',
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
    setValue('location', '');
    setValue('details', '');
    setValue('amount', 0);
    setValue('category', 'unspecified.unspecified');
  };

  const locationInputProps = register('location', {
    required: 'Please input a location',
  });
  const amountInputProps = register('amount', {
    min: { value: 1, message: 'Please input a positive amount' },
    required: 'Please input an amount',
  });
  const dateInputProps = register('date', { required: 'Please input a date' });
  const categoryInputProps = register('category');
  const detailsInputProps = register('details');

  return (
    <TxnFormBase
      title="Create Transaction"
      locationInputProps={locationInputProps}
      amountInputProps={amountInputProps}
      dateInputProps={dateInputProps}
      categoryInputProps={categoryInputProps}
      detailsInputProps={detailsInputProps}
      onSubmit={handleSubmit(submitCallback)}
    >
      <div className="mt-4 flex justify-end">
        <button className="rounded bg-indigo-500 py-2 px-4 font-bold text-white hover:bg-indigo-700 focus:outline-none focus:ring">
          Create transaction
        </button>
      </div>
    </TxnFormBase>
  );
}
