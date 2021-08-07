import { FormEvent, useEffect, useRef, useState } from "react";

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
  const [{ status, error }, setStatus] = useState({
    status: "idle",
    error: null,
  });
  const [statusMessage, setStatusMessage] = useState<string>();
  const cancelled = useRef(false);

  async function createUser(username: string, name: string) {
    const url = `${process.env.API_BASE_URL}/users`;
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, name }),
      });
      if (response.ok) {
        setStatusMessage(`User ${username} successfully created`);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function fetchUsers() {
    const url = `${process.env.API_BASE_URL}/users`;
    try {
      const response = await fetch(url);
      const parsed = await response.json();
      if (!cancelled.current) {
        setUsers(parsed.users);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setStatus({ status: "loading", error: null });
    const username = e.target.elements.username.value;
    const name = e.target.elements.name.value;
    try {
      await createUser(username, name);
      setStatus({ status: "fulfilled", error: null });
      await fetchUsers();
    } catch (err) {
      setStatus({ status: "rejected", error: err });
    }
  }

  useEffect(() => {
    fetchUsers();
    return () => {
      cancelled.current = true;
    };
  }, []);

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
      <form onSubmit={handleSubmit}>
        <span>
          <label htmlFor="name">Name</label>
          <input id="name" name="name" type="text" />
        </span>
        <span>
          <label htmlFor="username">Username</label>
          <input id="username" name="username" type="text" />
        </span>
        <span>
          <button type="submit">Create user</button>
        </span>
      </form>
      {status === "loading" && <p role="status">{status}</p>}
      {status === "fulfilled" && <p role="status">{statusMessage}</p>}
      {status === "rejected" && <p role="status">{error.message}</p>}
    </section>
  );
}
