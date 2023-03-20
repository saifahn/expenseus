import SharedLayout from 'components/LayoutShared';
import TrackersSubmitForm from 'components/TrackersSubmitForm';
import { useUserContext } from 'context/user';
import Link from 'next/link';
import useSWR from 'swr';

export interface Tracker {
  name: string;
  users: string[];
  id: string;
}

export default function SharedTrackers() {
  const { user } = useUserContext();
  const { data, error } = useSWR<Tracker[]>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/user/${user.id}`
      : null,
  );

  return (
    <SharedLayout>
      <TrackersSubmitForm />
      <div className="mt-4 p-4">
        <h2 className="text-2xl">My Trackers</h2>
        {error && <div>Failed to load: {error}</div>}
        {data === null && <div>Loading...</div>}
        {data && data.length === 0 && (
          <div>Please create a tracker to start</div>
        )}
        {data &&
          data.map((tracker) => (
            <Link href={`/shared/trackers/${tracker.id}`} key={tracker.id}>
              <a>
                <article className="mt-4 border-2 p-2">
                  <h3>{tracker.name}</h3>
                  {tracker.users.map((user) => (
                    <p key={user}>{user}</p>
                  ))}
                </article>
              </a>
            </Link>
          ))}
      </div>
    </SharedLayout>
  );
}
