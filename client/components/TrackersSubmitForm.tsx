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
      {allUsers && (
        <form onSubmit={handleSubmit(submitCallback)} className="border-4 p-6">
          <h3 className="text-lg font-semibold">Create Tracker</h3>
          <div className="mt-4">
            <label className="block font-semibold" htmlFor="name">
              Name
            </label>
            <input
              {...register('name', {
                required: 'Please input a tracker name',
              })}
              className="mt-2 w-full appearance-none rounded border py-2 px-3 leading-tight focus:outline-none focus:ring"
              type="text"
              id="name"
            />
          </div>
          <div className="mt-4">
            <label className="block font-semibold">Participants</label>
            <div>
              <p>
                {user.id} <span>(you)</span>
              </p>
            </div>
            <select
              {...register('partner', {
                required: 'A partner is required to create a tracker',
              })}
              className="mt-2 block rounded bg-white bg-clip-padding bg-no-repeat px-3 py-2 text-base font-normal text-gray-700 outline outline-1 transition ease-in-out focus:border-indigo-600 focus:bg-white focus:text-gray-700"
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
            <button className="rounded border-4 py-2 px-4 focus:outline-none focus:ring">
              Create tracker
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
