import AnalysisFormBase from 'components/AnalysisFormBase';
import { BarChart } from 'components/BarChart';
import PersonalLayout from 'components/LayoutPersonal';
import { fetcher } from 'config/fetcher';
import { useUserContext } from 'context/user';
import { subcategories, mainCategories } from 'data/categories';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { useSWRConfig } from 'swr';
import { Transaction } from 'types/Transaction';
import {
  calculateTotal,
  totalsByMainCategory,
  totalsBySubCategory,
} from 'utils/analysis';
import { plainDateStringToEpochSec, presets } from 'utils/dates';

type Inputs = {
  from: string;
  to: string;
};

export default function PersonalAnalysis() {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { register, handleSubmit, getValues, setValue } = useForm<Inputs>({
    shouldUseNativeValidation: true,
    defaultValues: {
      from: presets.ninetyDaysAgo().toString(),
      to: presets.now().toString(),
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
    <PersonalLayout>
      <AnalysisFormBase
        register={register}
        onSubmit={handleSubmit(submitCallback)}
        setValue={setValue}
      />
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
    </PersonalLayout>
  );
}
