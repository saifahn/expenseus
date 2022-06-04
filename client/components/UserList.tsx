import useSWR from 'swr';

export interface User {
  username: string;
  name: string;
  id: string;
}

export default function UserList() {
  const { data: users, error } = useSWR<User[]>(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`,
  );

  return (
    <section className="p-6 border-dotted border-4 border-indigo-800">
      <h2 className="text-2xl">Users</h2>
      {error && <div>Failed to load users</div>}
      {!error && !users && <div>Loading user information...</div>}
      {users &&
        users.map((user) => {
          return (
            <article
              className="p-4 mt-4 border-dotted border-2 border-pink-800"
              key={user.id}
            >
              <h3 className="text-xl">{user.username}</h3>
              <p>{user.name}</p>
              <p>{user.id}</p>
            </article>
          );
        })}
    </section>
  );
}
