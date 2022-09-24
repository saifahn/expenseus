import TrackerLayout from 'components/LayoutTracker';
import { fetcher } from 'config/fetcher';
import { mainCategories, subcategories } from 'data/categories';
import { useRouter } from 'next/router';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { useSWRConfig } from 'swr';
import {
  calculateTotal,
  totalsByMainCategory,
  totalsBySubCategory,
} from 'utils/analysis';
import { dateRanges, plainDateStringToEpochSec, presets } from 'utils/dates';
import { BarChart } from 'components/BarChart';
import { SharedTxn } from '.';
import { ChangeEvent } from 'react';

type Inputs = {
  from: string;
  to: string;
};

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

  function handlePresetSelect(e: ChangeEvent<HTMLSelectElement>) {
    const preset = e.target.value;
    const { from, to } = dateRanges[preset].presetFn();
    setValue('from', from);
    setValue('to', to);
  }

  return (
    <TrackerLayout>
      <form className="mt-4" onSubmit={handleSubmit(submitCallback)}>
        <h3 className="text-lg font-bold lowercase">Analyze transactions</h3>
        <div className="mt-3">
          <label className="block font-semibold lowercase text-slate-600">
            date preset
          </label>
          <select
            className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center lowercase placeholder-slate-400 focus:ring-0"
            onChange={handlePresetSelect}
          >
            {Object.entries(dateRanges).map(([preset, { name }]) => (
              <option key={preset} value={preset}>
                {name}
              </option>
            ))}
          </select>
        </div>
        <div className="mt-5">
          <label
            className="block font-semibold lowercase text-slate-600"
            htmlFor="dateFrom"
          >
            From
          </label>
          <input
            {...register('from', { required: 'Please input a date' })}
            className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center placeholder-slate-400 focus:ring-0"
            type="date"
            id="dateFrom"
          />
        </div>
        <div className="mt-5">
          <label
            className="block font-semibold lowercase text-slate-600"
            htmlFor="dateTo"
          >
            To
          </label>
          <input
            {...register('to', { required: 'Please input a date' })}
            className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center placeholder-slate-400 focus:ring-0"
            type="date"
            id="dateTo"
          />
        </div>

        <div className="mt-4 flex justify-end">
          <button
            className="rounded bg-violet-500 py-2 px-4 font-medium lowercase text-white hover:bg-violet-700 focus:outline-none focus:ring"
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
            <div className="h-screen">{BarChart(txns)}</div>
          </div>
        )}
      </div>
    </TrackerLayout>
  );
}
