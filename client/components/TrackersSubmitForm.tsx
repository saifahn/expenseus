import { useUserContext } from 'context/user';
import { SubmitHandler, useForm } from 'react-hook-form';
import useSWR, { useSWRConfig } from 'swr';
import { User } from './UserList';

type Inputs = {
  name: string;
  userId: string;
  partner: string;
};

async function createTracker(data: Inputs) {
  const payload = {
    name: data.name,
    users: [data.userId, data.partner],
  };
  await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers`, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(payload),
  });
}

export default function TrackersSubmitForm() {
  const { user } = useUserContext();
  const { mutate } = useSWRConfig();
  const { data: allUsers, error } = useSWR<User[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`,
  );

  const { register, handleSubmit, reset } = useForm<Inputs>({
    shouldUseNativeValidation: true,
  });

  const submitCallback: SubmitHandler<Inputs> = (data) => {
    if (!user) {
      console.error('user is not loaded');
      return;
    }
    data.userId = user.id;
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/user/${user.id}`,
      createTracker(data),
    );
    reset();
  };

  return (
    <div>
      {error && <div>Failed to load</div>}
      {!error && !allUsers && <div>Loading all users information...</div>}
      {allUsers && user && (
        <form onSubmit={handleSubmit(submitCallback)} className="bg-white">
          <h3 className="text-lg font-semibold lowercase">Create Tracker</h3>
          <div className="mt-4">
            <label className="block font-semibold lowercase" htmlFor="name">
              Name
            </label>
            <input
              {...register('name', {
                required: 'Please input a tracker name',
              })}
              className="focus:border-violet block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center text-xl placeholder-slate-400 focus:ring-0"
              type="text"
              id="name"
            />
          </div>
          <div className="mt-4">
            <label className="block font-semibold lowercase">
              Participants
            </label>
            <div className="mt-4 text-center">
              <p>
                {user.id} <span className="italic">(you)</span>
              </p>
              <p className="mt-2 font-medium">and</p>
            </div>
            <select
              {...register('partner', {
                required: 'A partner is required to create a tracker',
              })}
              className="focus:border-violet mt-2 block w-full appearance-none border-0 border-b-2 border-slate-200 px-4 text-center lowercase placeholder-slate-400 focus:ring-0"
            >
              {allUsers
                .filter((u) => u.id !== user.id)
                .map((u) => (
                  <option key={u.id} value={u.id}>
                    {u.id}
                  </option>
                ))}
            </select>
          </div>
          <div className="mt-6 flex justify-end">
            <button className="rounded bg-violet-500 py-2 px-4 font-medium lowercase text-white hover:bg-violet-700 focus:outline-none focus:ring">
              Create
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
