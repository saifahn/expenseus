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
import { plainDateStringToEpochSec, presets } from 'utils/dates';
import { BarChart } from 'components/BarChart';
import { SharedTxn } from '.';
import AnalysisFormBase from 'components/AnalysisFormBase';

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

  return (
    <TrackerLayout>
      <AnalysisFormBase
        register={register}
        onSubmit={handleSubmit(submitCallback)}
        setValue={setValue}
      />
      <div className="my-6">
        {error && <div>Failed to load details</div>}
        {txns === null && <div>Loading</div>}
        {txns?.length === 0 && <div>No transactions for that time period</div>}
        {txns?.length > 0 && (
          <div>
            <div className="h-screen">{BarChart(txns)}</div>
            <div className="mt-4">
              <p>
                In the period between {getValues().from} and {getValues().to},
                you have{' '}
                <span className="font-semibold">
                  {txns.length} transactions
                </span>
                , with a total cost of{' '}
                <span className="font-semibold">{calculateTotal(txns)}</span>.
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
          </div>
        )}
      </div>
    </TrackerLayout>
  );
}
