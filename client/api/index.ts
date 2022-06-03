import { User } from 'components/UserList';

const apiURL = process.env.NEXT_PUBLIC_API_BASE_URL;

export class UserAPI {
  baseURL = `${apiURL}/users`;

  async listUsers() {
    const res = await fetch(this.baseURL, {
      credentials: 'include',
    });
    const parsed: User[] = await res.json();
    return parsed;
  }

  async getSelf() {
    const res = await fetch(`${this.baseURL}/self`, {
      credentials: 'include',
    });
    if (!res.ok) {
      throw res.status;
    }
    const parsed: User = await res.json();
    return parsed;
  }
}

export const fetcher = (url) =>
  fetch(url, { credentials: 'include' }).then((res) => res.json());
