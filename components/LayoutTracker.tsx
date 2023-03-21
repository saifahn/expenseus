import Link from 'next/link';
import { useRouter } from 'next/router';
import { Tracker } from 'ddb/trackers';
import { PropsWithChildren } from 'react';
import useSWR from 'swr';
import SharedLayout from './LayoutShared';

export default function TrackerLayout({ children }: PropsWithChildren<{}>) {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data: tracker, error } = useSWR<Tracker>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}`,
  );

  return (
    <SharedLayout>
      {error && <div>Failed to load</div>}
      {!error && !tracker && <div>Loading tracker information...</div>}
      {tracker && (
        <>
          <Link href={`/shared/trackers/${trackerId}`}>
            <a>
              <h2 className="mt-4 text-2xl underline decoration-rose-300 decoration-2 hover:decoration-rose-500">
                {tracker.name}
              </h2>
            </a>
          </Link>
          {children}
        </>
      )}
    </SharedLayout>
  );
}
