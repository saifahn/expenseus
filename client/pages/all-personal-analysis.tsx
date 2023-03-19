import AnalysisFormBase from 'components/AnalysisFormBase';
import { BarChart } from 'components/BarChart';
import { fetcher } from 'config/fetcher';
import { useUserContext } from 'context/user';
import { mainCategories, subcategories } from 'data/categories';
import { AllTxnsResponse } from 'pages';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { mutate } from 'swr';
import { Transaction } from 'types/Transaction';
import {
  calculatePersonalTotal,
  calculateTotal,
  totalsByMainCategory,
  totalsBySubCategory,
} from 'utils/analysis';
import { plainDateStringToEpochSec, presets } from 'utils/dates';
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
      (a, b) => b.date - a.date,
    );
  }
  let total;
  if (user) {
    total = calculatePersonalTotal(user.id, txns);
  }

  const submitCallback: SubmitHandler<Inputs> = () => {
    mutate('all.analysis');
  };

  return (
    <>
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
        {txns && txns.length > 0 && (
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
    </>
  );
}
