import SharedLayout from 'components/LayoutShared';
import { useUserContext } from 'context/user';
import Link from 'next/link';
import useSWR from 'swr';
import { Tracker } from './trackers';

export default function SharedIndex() {
  const { user } = useUserContext();
  const { data, error } = useSWR<Tracker[]>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/user/${user.id}`
      : null,
  );

  return (
    <SharedLayout>
      {error && <div>Failed to load: {error}</div>}
      {data === null && <div>Loading...</div>}
      {data && data.length === 0 && <div>You are not part of any trackers</div>}
      {data &&
        data.map((tracker) => (
          <Link href={`/shared/trackers/${tracker.id}`} key={tracker.id}>
            <a>
              <article className="mt-3 rounded-lg border-2 border-slate-200 p-3 hover:bg-slate-200 active:bg-slate-300">
                <p className="text-xl font-semibold">{tracker.name}</p>
                <div className="mt-2">
                  {tracker.users.map((user) => (
                    <p key={user}>{user}</p>
                  ))}
                </div>
              </article>
            </a>
          </Link>
        ))}
      <Link href="/shared/trackers/create">
        <a className="mt-4 block rounded-lg bg-violet-50 p-3 font-medium lowercase text-black hover:bg-violet-100 active:bg-violet-200">
          Create new tracker +
        </a>
      </Link>
    </SharedLayout>
  );
}
