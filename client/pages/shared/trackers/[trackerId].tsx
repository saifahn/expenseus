import { fetcher } from 'api';
import SharedLayout from 'components/SharedLayout';
import { useRouter } from 'next/router';
import useSWR from 'swr';
import { Tracker } from '.';

export default function TrackerPage() {
  const router = useRouter();
  const { trackerId } = router.query;
  const { data, error } = useSWR<Tracker>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/${trackerId}`,
    fetcher,
  );

  return (
    <SharedLayout>
      {error && <div>Failed to load</div>}
      {!error && !data && <div>Loading...</div>}
      {data && (
        <>
          <h2 className="text-2xl mt-8">{data.name}</h2>
          <h3 className="mt-2">{data.id}</h3>
          {data.users.map((user) => (
            <p key={user}>{user}</p>
          ))}
        </>
      )}
    </SharedLayout>
  );
}
