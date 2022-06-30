import PersonalLayout from 'components/LayoutPersonal';
import { fetcher } from 'config/fetcher';
import { useUserContext } from 'context/user';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { useSWRConfig } from 'swr';
import { Temporal } from 'temporal-polyfill';
import { Transaction } from 'types/Transaction';
import { plainDateStringToEpochSec } from 'utils/temporal';

type Inputs = {
  from: string;
  to: string;
};

function calculateTotal(txns: Transaction[]) {
  let total = 0;
  for (const txn of txns) {
    total += txn.amount;
  }
  return total;
}

export default function PersonalAnalysis() {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, getValues } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      from: Temporal.Now.plainDateISO().subtract({ months: 3 }).toString(),
      to: Temporal.Now.plainDateISO().toString(),
    },
  });

  const { data: txns, error } = useSWR<Transaction[]>(
    'personal.analysis',
    () => {
      if (!user) return null;
      const { from, to } = getValues();
      const fromEpochSec = plainDateStringToEpochSec(from);
      const toEpochSec = plainDateStringToEpochSec(to);
      return fetcher(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}/range?from=${fromEpochSec}&to=${toEpochSec}`,
      );
    },
  );

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate('personal.analysis');
  };

  return (
    <>
      <PersonalLayout>
        <form className="border-4 p-6" onSubmit={handleSubmit(submitCallback)}>
          <div className="mt-4">
            <label className="block font-semibold" htmlFor="dateFrom">
              From
            </label>
            <input
              {...register('from', { required: 'Please input a date' })}
              className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
              type="date"
              id="dateFrom"
            />
          </div>
          <div className="mt-4">
            <label className="block font-semibold" htmlFor="dateTo">
              To
            </label>
            <input
              {...register('to', { required: 'Please input a date' })}
              className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
              type="date"
              id="dateTo"
            />
          </div>
          <div className="mt-4 flex justify-end">
            <button className="rounded bg-indigo-500 py-2 px-4 text-sm font-bold uppercase text-white hover:bg-indigo-700 focus:outline-none focus:ring">
              Get details
            </button>
          </div>
        </form>
        <div className="mt-6">
          {error && <div>Failed to load details</div>}
          {txns === null && <div>Loading</div>}
          {txns?.length === 0 && (
            <div>No transactions for that time period</div>
          )}
          {txns?.length > 0 && (
            <div>
              <h3>In that time period:</h3>
              <ul>
                <li>You have {txns.length} transaction(s)</li>
                <li>You have spent {calculateTotal(txns)}</li>
              </ul>
            </div>
          )}
        </div>
      </PersonalLayout>
    </>
  );
}
