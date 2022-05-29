import SharedLayout from 'components/SharedLayout';
import TrackersSubmitForm from 'components/TrackersSubmitForm';
import { useUserContext } from 'context/user';
import useSWR from 'swr';

interface Tracker {
  name: string;
  users: string[];
  id: string;
}

const fetcher = (url) =>
  fetch(url, { credentials: 'include' }).then((res) => res.json());

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
        {data &&
          data.map((tracker) => (
            <article className="p-2 border-2 mt-4" key={tracker.id}>
              <h3>{tracker.name}</h3>
              {tracker.users.map((user) => (
                <p key={user}>{user}</p>
              ))}
            </article>
          ))}
      </div>
    </SharedLayout>
  );
}
