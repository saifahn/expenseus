import { Transaction } from "components/TransactionList";
import { User } from "components/UserList";
import { v4 as uuidv4 } from "uuid";

const apiURL = process.env.NEXT_PUBLIC_API_BASE_URL;

export class UserAPI {
  baseURL = `${apiURL}/users`;

  async listUsers() {
    const res = await fetch(this.baseURL, {
      credentials: "include",
    });
    const parsed: User[] = await res.json();
    return parsed;
  }

  async createUser(username: string, name: string) {
    // TODO: remove this once the back end handles id creation
    const id = uuidv4();
    const res = await fetch(this.baseURL, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username, name, id }),
    });
    if (res.ok) {
      // TODO: make this depend on the server response?
      return `User ${username} was successfully created`;
    }
    // TODO: handle this better?
    throw new Error(res.statusText);
  }

  async getSelf() {
    const res = await fetch(`${this.baseURL}/self`, {
      credentials: "include",
    });
    if (!res.ok) {
      throw res.status;
    }
    const parsed: User = await res.json();
    return parsed;
  }
}

export class TransactionAPI {
  baseURL = `${apiURL}/transactions`;

  async listTransactions() {
    const res = await fetch(this.baseURL, {
      credentials: "include",
    });
    const parsed: Transaction[] = await res.json();
    return parsed;
  }

  async createTransaction(data: FormData) {
    const res = await fetch(this.baseURL, {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      credentials: "include",
      body: data,
    });
    if (res.ok) {
      // TODO: return the message from the server
      return `Transaction was successfully created`;
    }
    // TODO: handle this better?
    throw new Error(res.statusText);
  }

  async getTransaction(id: string) {
    const res = await fetch(`${this.baseURL}/${id}`, {
      credentials: "include",
    });
    const parsed: Transaction = await res.json();
    return parsed;
  }
}
