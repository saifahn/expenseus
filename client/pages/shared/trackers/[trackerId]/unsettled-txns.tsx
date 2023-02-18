import TrackerLayout from 'components/LayoutTracker';
import { useRouter } from 'next/router';
import useSWR, { mutate } from 'swr';
import { SharedTxn, SharedTxnOne } from '.';

type UnsettledResponse = {
  transactions: SharedTxn[];
  debtor: string;
  debtee: string;
  amountOwed: number;
};

async function settleUnsettledTxns(txns: SharedTxn[]) {
  await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/shared/settle`,
    {
      method: 'POST',
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(txns),
    },
  );
}

export default function UnsettledTxnPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data: response, error } = useSWR<UnsettledResponse>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions/unsettled`,
  );

  function handleSettleUp() {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}/transactions/unsettled`,
      // will only be called when there is a response
      settleUnsettledTxns(response!.transactions),
    );
  }

  return (
    <TrackerLayout>
      {error && <div>Error loading unsettled transactions</div>}
      {response === null && <div>Loading unsettled transactions</div>}
      {response?.transactions === null && (
        <p className="mt-4">You currently have no unsettled transactions!</p>
      )}
      {response && response.transactions?.length > 0 && (
        <>
          <p className="mt-4">
            {response.debtor} owes {response.debtee} {response.amountOwed} for{' '}
            {response.transactions.length} transactions
          </p>
          {response.transactions.map((txn) => (
            <SharedTxnOne txn={txn} key={txn.id} />
          ))}
          <button
            className="mt-4 rounded bg-violet-100 py-2 px-4 font-medium lowercase hover:bg-violet-200 focus:outline-none focus:ring active:bg-violet-300"
            onClick={handleSettleUp}
          >
            Settle up
          </button>
        </>
      )}
    </TrackerLayout>
  );
}
