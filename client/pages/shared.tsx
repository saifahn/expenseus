import TrackersSubmitForm from 'components/TrackersSubmitForm';
import { useUserContext } from 'context/user';
import { useState } from 'react';
import useSWR from 'swr';

interface Tracker {
  name: string;
  users: string[];
  id: string;
}

const fetcher = (url) =>
  fetch(url, { credentials: 'include' }).then((res) => res.json());

export default function Shared() {
  const { user } = useUserContext();
  const [showing, setShowing] = useState<'home' | 'trackers'>('home');
  const { data, error } = useSWR<Tracker[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/user/${user.id}`,
    fetcher,
  );

  return (
    <>
      <h1 className="text-4xl">Shared</h1>
      <nav className="mt-4">
        <ul className="flex">
          <li
            onClick={() => setShowing('home')}
            className="p-2 border-2 cursor-pointer"
          >
            Home
          </li>
          <li
            onClick={() => setShowing('trackers')}
            className="p-2 border-2 cursor-pointer ml-4"
          >
            Trackers
          </li>
        </ul>
      </nav>
      <section className="mt-4">
        {showing === 'home' && <p>Showing list of shared transactions</p>}
        {showing === 'trackers' && (
          <>
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
          </>
        )}
      </section>
    </>
  );
}
