import AnalysisFormBase from 'components/AnalysisFormBase';
import { BarChart, personalTotalsForBarChart } from 'components/BarChart';
import { fetcher } from 'config/fetcher';
import { useUserContext } from 'context/user';
import { mainCategories, subcategories } from 'data/categories';
import Head from 'next/head';
import { AllTxnsResponse } from 'pages';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { mutate } from 'swr';
import { Transaction } from 'types/Transaction';
import {
  calculatePersonalTotal,
  personalTotalsByCategory,
} from 'utils/analysis';
import { plainDateStringToEpochSec, presets } from 'utils/dates';
import { jpyFormatter } from 'utils/jpyFormatter';
import { SharedTxn } from './shared/trackers/[trackerId]';

type Inputs = {
  from: string;
  to: string;
};

/**
 * This represents the page for analyzing all transactions that the logged-in
 * user is involved in
 */
export default function AllAnalysis() {
  const { user } = useUserContext();
  const { register, handleSubmit, getValues, setValue } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      from: presets.ninetyDaysAgo().toString(),
      to: presets.now().toString(),
    },
  });

  const { data: allTxns, error } = useSWR<AllTxnsResponse>(
    'all.analysis',
    async () => {
      if (!user) return null;
      const { from, to } = getValues();
      const fromEpochSec = plainDateStringToEpochSec(from);
      const toEpochSec = plainDateStringToEpochSec(to);
      return fetcher(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}/all?from=${fromEpochSec}&to=${toEpochSec}`,
      );
    },
  );

  let txns: (Transaction | SharedTxn)[] = [];
  if (allTxns) {
    txns = [...allTxns.transactions, ...allTxns.sharedTransactions].sort(
      (a, b) => a.date - b.date, // sorted ascending for the graph so it displays most recent last
    );
  }

  const submitCallback: SubmitHandler<Inputs> = () => {
    mutate('all.analysis');
  };

  return (
    <>
      <Head>
        <title>all personal analysis - expenseus</title>
      </Head>
      <h1 className="mb-2 text-xl font-semibold lowercase">
        All personal spending analysis
      </h1>
      <AnalysisFormBase
        register={register}
        onSubmit={handleSubmit(submitCallback)}
        setValue={setValue}
      />
      <div className="my-6">
        {error && <div>Failed to load details</div>}
        {txns === null && <div>Loading</div>}
        {txns?.length === 0 && <div>No transactions for that time period</div>}
        {user && txns?.length > 0 && (
          <div>
            <div className="h-screen">
              {BarChart(personalTotalsForBarChart(txns, user.id))}
            </div>
            <div className="mt-4">
              <p>
                In the period between {getValues().from} and {getValues().to},
                you have{' '}
                <span className="font-semibold">
                  {txns.length} transactions
                </span>
                , with a total cost of{' '}
                <span className="font-semibold">
                  {jpyFormatter.format(calculatePersonalTotal(user.id, txns))}
                </span>
                .
              </p>
              <p className="mt-4 text-lg font-medium">Main categories:</p>
              {personalTotalsByCategory(user.id, txns).map((cat) => (
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
    </>
  );
}
