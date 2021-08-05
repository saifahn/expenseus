import { useEffect, useState } from "react";

export interface User {
  username: string;
  name: string;
  id: string;
}

/**
 * Component rendering a list of users
 */
export default function Users() {
  const [users, setUsers] = useState<User[]>();
  useEffect(() => {
    let cancelled = false;

    async function fetchUsers() {
      let url = `${process.env.API_BASE_URL}/users`;
      try {
        const response = await fetch(url);
        const parsed = await response.json();
        if (!cancelled) {
          setUsers(parsed.users);
        }
      } catch (err) {
        console.error(err);
      }
    }

    fetchUsers();
    return () => {
      cancelled = true;
    };
  });
  return (
    <section>
      <h2>Users</h2>
      {users &&
        users.map(user => {
          return (
            <article key={user.id}>
              <h3>{user.username}</h3>
              <p>{user.name}</p>
              <p>{user.id}</p>
            </article>
          );
        })}
      {/* <form onSubmit={handleSubmit}></form> */}
    </section>
  );
}
