import { useUserContext } from 'context/user';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useSWRConfig } from 'swr';

type Inputs = {
  name: string;
  users: string;
};

async function createTracker(data: Inputs) {
  const payload = {
    name: data.name,
    users: data.users.split(','),
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
  const { register, handleSubmit } = useForm<Inputs>({
    defaultValues: {
      users: user.id,
    },
  });
  const submitCallback: SubmitHandler<Inputs> = (data) => {
    mutate(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/trackers/user/${user.id}`,
      createTracker(data),
    );
  };

  return (
    <div className="mt-6">
      <form onSubmit={handleSubmit(submitCallback)} className="border-4 p-6">
        <h3 className="text-lg font-semibold">Create Tracker</h3>
        <div className="mt-4">
          <label className="block font-semibold" htmlFor="name">
            Name
          </label>
          <input
            {...register('name')}
            className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
            type="text"
            id="name"
          />
        </div>
        {/* TODO: make this a list of users that is populated by get all users */}
        <div className="mt-4">
          <label className="block font-semibold" htmlFor="participants">
            Participants
          </label>
          <input
            {...register('users')}
            className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
            type="text"
            id="participants"
          />
        </div>
        <div className="mt-6 flex justify-end">
          <button className="border-4 py-2 px-4 rounded focus:outline-none focus:ring">
            Create tracker
          </button>
        </div>
      </form>
    </div>
  );
}
