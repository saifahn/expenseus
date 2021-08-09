import { Expense } from "components/ExpenseList";
import { User } from "components/UserList";
import { v4 as uuidv4 } from "uuid";

const apiURL = process.env.NEXT_PUBLIC_API_BASE_URL;

export class UserAPI {
  baseURL = `${apiURL}/users`;

  async listUsers() {
    const url = `${this.baseURL}`;
    const res = await fetch(url);
    const parsed: User[] = await res.json();
    return parsed;
  }

  async createUser(username: string, name: string) {
    // TODO: remove this once the back end handles id creation
    const id = uuidv4();
    const url = `${this.baseURL}`;
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, name, id }),
    });
    if (res.ok) {
      // TODO: make this depend on the server response?
      return `User ${username} was successfully created`;
    }
  }
}

export class ExpenseAPI {
  baseURL = `${apiURL}/expenses`;

  async listExpenses() {
    const res = await fetch(this.baseURL);
    const parsed: Expense[] = await res.json();
    return parsed;
  }
}
