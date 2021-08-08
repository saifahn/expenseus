import { FormEvent, useEffect, useRef, useState } from "react";
import { v4 as uuidv4 } from "uuid";

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
  const [newUsername, setNewUsername] = useState("");
  const [newName, setNewName] = useState("");
  const cancelled = useRef(false);

  async function createUser(username: string, name: string) {
    const url = `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`;
    // TODO: remove this once the back end handles id creation
    const id = uuidv4();
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, name, id }),
      });
      if (response.ok) {
        setStatusMessage(`User ${username} successfully created`);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function fetchUsers() {
    const url = `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`;
    try {
      const response = await fetch(url);
      const parsed = await response.json();
      if (!cancelled.current) {
        setUsers(parsed);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setStatus({ status: "loading", error: null });
    try {
      await createUser(newUsername, newName);
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
          <input
            id="name"
            name="name"
            type="text"
            value={newName}
            onChange={e => setNewName(e.target.value)}
          />
        </span>
        <span>
          <label htmlFor="username">Username</label>
          <input
            id="username"
            name="username"
            type="text"
            value={newUsername}
            onChange={e => setNewUsername(e.target.value)}
          />
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
