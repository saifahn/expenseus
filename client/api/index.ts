import { User } from "components/UserList";
const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL;

async function listUsers() {
  const url = `${baseURL}/users`;
  const res = await fetch(url);
  const parsed: User[] = await res.json();
  return parsed;
}

export { listUsers };
