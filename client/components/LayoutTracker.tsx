import Link from 'next/link';
import { useRouter } from 'next/router';
import { Tracker } from 'pages/shared/trackers';
import useSWR from 'swr';
import SharedLayout from './LayoutShared';

export default function TrackerLayout({ children }) {
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
          <h2 className="mt-8 text-2xl">{tracker.name}</h2>
          <h3 className="mt-2">{tracker.id}</h3>
          {tracker.users.map((user) => (
            <p key={user}>{user}</p>
          ))}
          <nav className="mt-4">
            <ul className="flex">
              <li className="flex">
                <Link href={`/shared/trackers/${trackerId}`}>
                  <a className="border-2 p-2">Transactions</a>
                </Link>
              </li>
              <li className="flex">
                <Link href={`/shared/trackers/${trackerId}/unsettled-txns`}>
                  <a className="ml-4 border-2 p-2">Unsettled</a>
                </Link>
              </li>
              <li className="flex">
                <Link href={`/shared/trackers/${trackerId}/analysis`}>
                  <a className="ml-4 border-2 p-2">Analysis</a>
                </Link>
              </li>
            </ul>
          </nav>
          {children}
        </>
      )}
    </SharedLayout>
  );
}
