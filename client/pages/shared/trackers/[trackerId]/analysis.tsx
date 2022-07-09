import TrackerLayout from 'components/LayoutTracker';
import { fetcher } from 'config/fetcher';
import {
  mainCategories,
  MainCategoryKey,
  subcategories,
  SubcategoryKey,
} from 'data/categories';
import { useRouter } from 'next/router';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { useSWRConfig } from 'swr';
import { dateRanges, plainDateStringToEpochSec, presets } from 'utils/dates';
import { SharedTxn } from '.';

type Inputs = {
  from: string;
  to: string;
};

function calculateTotal(txns: SharedTxn[]) {
  let total = 0;
  for (const txn of txns) {
    total += txn.amount;
  }
  return total;
}

function totalsByMainCategory(txns: SharedTxn[]) {
  const totals = {} as Record<MainCategoryKey, number>;
  for (const txn of txns) {
    const mainCategory = subcategories[txn.category].mainCategory;
    if (!totals[mainCategory]) totals[mainCategory] = 0;
    totals[mainCategory] += txn.amount;
  }
  return Object.entries(totals).map(([category, total]) => ({
    category,
    total,
  }));
}

function totalsBySubCategory(txns: SharedTxn[]) {
  const totals = {} as Record<SubcategoryKey, number>;
  for (const txn of txns) {
    if (!totals[txn.category]) totals[txn.category] = 0;
    totals[txn.category] += txn.amount;
  }
  return Object.entries(totals).map(([category, total]) => ({
    category,
    total,
  }));
}

export default function TrackerAnalysis() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, getValues, setValue } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      from: presets.ninetyDaysAgo().toString(),
      to: presets.now().toString(),
    },
  });

  const { data: txns, error } = useSWR<SharedTxn[]>(
    `${trackerId}.analysis`,
    () => {
      if (!trackerId) return null;
      const { from, to } = getValues();
      const fromEpochSec = plainDateStringToEpochSec(from);
      const toEpochSec = plainDateStringToEpochSec(to);
      return fetcher(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions/range?from=${fromEpochSec}&to=${toEpochSec}`,
      );
    },
  );

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate(`${trackerId}.analysis`);
  };

  function handlePresetClick(presetFn) {
    const { from, to } = presetFn();
    setValue('from', from);
    setValue('to', to);
  }

  return (
    <TrackerLayout>
      <form
        className="mt-4 border-4 p-6"
        onSubmit={handleSubmit(submitCallback)}
      >
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
        <div className="mt-4">
          {Object.entries(dateRanges).map(([preset, { name, presetFn }]) => (
            <button
              onClick={() => handlePresetClick(presetFn)}
              key={preset}
              className="mr-2 rounded border-2 border-indigo-300 py-2 px-4 text-sm hover:border-indigo-700 focus:outline-none focus:ring"
            >
              {name}
            </button>
          ))}
        </div>
        <div className="mt-4 flex justify-end">
          <button
            className="rounded bg-indigo-500 py-2 px-4 text-sm font-bold uppercase text-white hover:bg-indigo-700 focus:outline-none focus:ring"
            type="submit"
          >
            Get details
          </button>
        </div>
      </form>
      <div className="mt-6">
        {error && <div>Failed to load details</div>}
        {txns === null && <div>Loading</div>}
        {txns?.length === 0 && <div>No transactions for that time period</div>}
        {txns?.length > 0 && (
          <div>
            <h3 className="text-xl font-medium">In this time period:</h3>
            <p>
              You have {txns.length} transactions, with a total cost of{' '}
              {calculateTotal(txns)}
            </p>
            <p className="mt-4 text-lg font-medium">Main categories:</p>
            <ul className="list-inside list-disc">
              {totalsByMainCategory(txns).map((total) => (
                <li key={total.category}>
                  You spent {total.total} on{' '}
                  {mainCategories[total.category].en_US}
                </li>
              ))}
            </ul>
            <p className="mt-4 text-lg font-medium">Subcategories:</p>
            <ul className="list-inside list-disc">
              {totalsBySubCategory(txns).map((total) => (
                <li key={total.category}>
                  You spent {total.total} on{' '}
                  {subcategories[total.category].en_US}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </TrackerLayout>
  );
}
