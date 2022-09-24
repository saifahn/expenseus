import { useUserContext } from 'context/user';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';
import { plainDateISONowString } from 'utils/dates';
import TxnFormBase, { createTxnFormData, TxnFormInputs } from './TxnFormBase';

async function createTransaction(data: TxnFormInputs) {
  const formData = createTxnFormData(data);

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
  const router = useRouter();
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, setValue } = useForm<TxnFormInputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      location: '',
      details: '',
      amount: null,
      date: plainDateISONowString(),
      category: 'unspecified.unspecified',
    },
  });

  const submitCallback: SubmitHandler<TxnFormInputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`,
      createTransaction(data),
    );
    setValue('location', '');
    setValue('details', '');
    setValue('amount', null);
    setValue('category', 'unspecified.unspecified');
  };

  return (
    <TxnFormBase
      title="Create Transaction"
      register={register}
      onSubmit={handleSubmit(submitCallback)}
    >
      <div className="mt-5 flex justify-end">
        <Link href="/personal">
          <a className="mr-2 rounded py-2 px-4 font-medium lowercase hover:bg-slate-200 focus:outline-none focus:ring">
            Close
          </a>
        </Link>
        <button className="rounded bg-violet-500 py-2 px-4 font-medium lowercase text-white hover:bg-violet-700 focus:outline-none focus:ring">
          Create
        </button>
      </div>
    </TxnFormBase>
  );
}
