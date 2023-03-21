import TrackerLayout from 'components/LayoutTracker';
import SharedTxnCreateForm from 'components/SharedTxnCreateForm';
import { useRouter } from 'next/router';
import useSWR from 'swr';
import { Tracker } from 'ddb/trackers';
import Head from 'next/head';

export default function CreateSharedTransaction() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data: tracker, error } = useSWR<Tracker>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}`,
  );
  return (
    <TrackerLayout>
      <Head>
        <title>create shared transaction - expenseus</title>
      </Head>
      {error && <div>Failed to load: {error}</div>}
      {tracker === null && <div>Loading...</div>}
      {tracker && (
        <div className="mt-4">
          <SharedTxnCreateForm tracker={tracker} />
        </div>
      )}
    </TrackerLayout>
  );
}
