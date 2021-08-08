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
      <div className="mt-6">
        <h2 className="text-2xl">Create users</h2>
        <div className="mx-auto w-full max-w-xs">
          <form
            className="bg-white p-6 rounded-lg shadow-md"
            onSubmit={handleSubmit}
          >
            <div>
              <label className="block font-semibold" htmlFor="name">
                Name
              </label>
              <input
                className="shadow appearance-none w-full border rounded mt-2 py-2 px-3 leading-tight focus:outline-none focus:ring"
                id="name"
                name="name"
                type="text"
                value={newName}
                onChange={e => setNewName(e.target.value)}
              />
            </div>
            <div className="mt-6">
              <label className="block font-semibold" htmlFor="username">
                Username
              </label>
              <input
                className="shadow appearance-none w-full border rounded mt-2 py-2 px-3 leading-tight focus:outline-none focus:ring"
                id="username"
                name="username"
                type="text"
                value={newUsername}
                onChange={e => setNewUsername(e.target.value)}
              />
            </div>
            <div className="mt-6 flex justify-end">
              <button
                className="bg-indigo-500 hover:bg-indigo-700 text-white py-2 px-4 rounded focus:outline-none focus:ring"
                type="submit"
              >
                Create user
              </button>
            </div>
          </form>
          {status === "loading" && <p role="status">{status}</p>}
          {status === "fulfilled" && <p role="status">{statusMessage}</p>}
          {status === "rejected" && <p role="status">{error.message}</p>}
        </div>
      </div>
    </section>
  );
}
