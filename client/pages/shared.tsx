import TrackersSubmitForm from 'components/TrackersSubmitForm';
import { useState } from 'react';

export default function Shared() {
  const [showing, setShowing] = useState<'home' | 'trackers'>('home');

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
          </>
        )}
      </section>
    </>
  );
}
