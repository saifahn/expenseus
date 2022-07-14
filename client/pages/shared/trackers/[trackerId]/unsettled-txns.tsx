import TrackerLayout from 'components/LayoutTracker';
import { useRouter } from 'next/router';
import useSWR from 'swr';
import { SharedTxn } from '.';

type UnsettledResponse = {
  transactions: SharedTxn[];
  debtor: string;
  debtee: string;
  amountOwed: number;
};

export default function UnsettledTxnPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data: response, error } = useSWR<UnsettledResponse>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions/unsettled`,
  );

  return (
    <TrackerLayout>
      {error && <div>Error loading unsettled transactions</div>}
      {response === null && <div>Loading unsettled transactions</div>}
      {response?.transactions === null && (
        <p className="mt-4">You currently have no unsettled transactions!</p>
      )}
      {response?.transactions?.length > 0 && (
        <p className="mt-4">
          {response.debtor} owes {response.debtee} {response.amountOwed} for{' '}
          {response.transactions.length} transactions
        </p>
      )}
    </TrackerLayout>
  );
}
