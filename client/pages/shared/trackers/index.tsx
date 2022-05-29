import { fetcher } from 'api';
import SharedLayout from 'components/SharedLayout';
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
  const { data, error } = useSWR<Tracker[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/user/${user.id}`,
    fetcher,
  );

  return (
    <SharedLayout>
      <TrackersSubmitForm />
      <div className="mt-4 p-4">
        <h2 className="text-2xl">My Trackers</h2>
        {error && <div>Failed to load: {error}</div>}
        {!data && <div>Loading...</div>}
        {data &&
          data.map((tracker) => (
            <Link href={`/shared/trackers/${tracker.id}`} key={tracker.id}>
              <a>
                <article className="p-2 border-2 mt-4">
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