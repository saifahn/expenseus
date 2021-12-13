import { useEffect, useRef, useState } from "react";
import { UserAPI } from "api";

export interface User {
  username: string;
  name: string;
  id: string;
}

/**
 * Component rendering a list of users
 */
export default function UserList() {
  const [users, setUsers] = useState<User[]>();
  const cancelled = useRef(false);

  async function fetchUsers() {
    try {
      const api = new UserAPI();
      const users = await api.listUsers();
      if (!cancelled.current) {
        setUsers(users);
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    fetchUsers();
    return () => {
      cancelled.current = true;
    };
  }, []);

  return (
    <section className="p-6 shadow-lg bg-indigo-50 rounded-xl">
      <h2 className="text-2xl">Users</h2>
      {users &&
        users.map(user => {
          return (
            <article
              className="p-4 mt-4 rounded-md shadow-md bg-white"
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
