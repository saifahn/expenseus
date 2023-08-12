import TrackerLayout from 'components/LayoutTracker';
import { fetcher } from 'config/fetcher';
import { mainCategories, subcategories } from 'data/categories';
import { useRouter } from 'next/router';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { useSWRConfig } from 'swr';
import { calculateTotal, totalsByCategory } from 'utils/analysis';
import { plainDateStringToEpochSec, presets } from 'utils/dates';
import { BarChart, totalsForBarChart } from 'components/BarChart';
import { SharedTxn } from '.';
import AnalysisFormBase from 'components/AnalysisFormBase';
import Head from 'next/head';
import { jpyFormatter } from 'utils/jpyFormatter';

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
    async () => {
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
      <Head>
        <title>analyze shared transactions - expenseus</title>
      </Head>
      <AnalysisFormBase
        register={register}
        onSubmit={handleSubmit(submitCallback)}
        setValue={setValue}
      />
      <div className="my-6">
        {error && <div>Failed to load details</div>}
        {txns === null && <div>Loading</div>}
        {txns?.length === 0 && <div>No transactions for that time period</div>}
        {txns && txns.length > 0 && (
          <div>
            <div className="h-screen">{BarChart(totalsForBarChart(txns))}</div>
            <div className="mt-4">
              <p>
                In the period between {getValues().from} and {getValues().to},
                you have{' '}
                <span className="font-semibold">
                  {txns.length} transactions
                </span>
                , with a total cost of{' '}
                <span className="font-semibold">
                  {jpyFormatter.format(calculateTotal(txns))}
                </span>
                .
              </p>
              {totalsByCategory(txns).map((cat) => (
                <div key={cat.mainCategory}>
                  <p className="mt-3 text-lg font-medium">
                    {mainCategories[cat.mainCategory].emoji}
                    &nbsp;
                    {mainCategories[cat.mainCategory].en_US}
                    {' - '}
                    {jpyFormatter.format(cat.total)}
                  </p>
                  <ul className="mt-1">
                    {cat.subcategories.map((subcat) => (
                      <li key={subcat.category}>
                        {subcategories[subcat.category].en_US}
                        {' - '}
                        {jpyFormatter.format(subcat.total)}
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </TrackerLayout>
  );
}
