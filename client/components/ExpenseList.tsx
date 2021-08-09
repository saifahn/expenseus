import { ExpenseAPI } from "api";
import { useState, useRef, useEffect, FormEvent } from "react";

export interface Expense {
  userID: string;
  name: string;
  id: string;
}

export default function ExpenseList() {
  const [expenses, setExpenses] = useState<Expense[]>();
  const [expenseName, setExpenseName] = useState("");
  const [expenseUserID, setExpenseUserID] = useState("");
  const [{ status, error }, setStatus] = useState({
    status: "idle",
    error: null,
  });
  const [statusMessage, setStatusMessage] = useState<string>();
  const cancelled = useRef(false);

  async function fetchExpenses() {
    const url = `${process.env.NEXT_PUBLIC_API_BASE_URL}/expenses`;
    try {
      const api = new ExpenseAPI();
      const expenses = await api.listExpenses();
      if (!cancelled.current) {
        setExpenses(expenses);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function createExpense(name: string, userID: string) {
    const url = `${process.env.NEXT_PUBLIC_API_BASE_URL}/expenses`;
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, userid: userID }),
      });
      if (response.ok) {
        setStatusMessage(`Expense ${name} successfully created`);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setStatus({ status: "loading", error: null });
    try {
      await createExpense(expenseName, expenseUserID);
      setStatus({ status: "fulfilled", error: null });
      await fetchExpenses();
    } catch (err) {
      setStatus({ status: "rejected", error: err });
    }
  }

  useEffect(() => {
    fetchExpenses();
    return () => {
      cancelled.current = true;
    };
  }, []);

  return (
    <section className="p-6 mt-8 shadow-lg bg-indigo-50 rounded-xl">
      <h2 className="text-2xl">Expenses</h2>
      {expenses &&
        expenses.map(expense => {
          return (
            <article
              className="p-4 mt-4 rounded-md shadow-md bg-white"
              key={expense.id}
            >
              <h3 className="text-xl">{expense.name}</h3>
              <p>{expense.userID}</p>
              <p>{expense.id}</p>
            </article>
          );
        })}
      <div className="mt-6">
        <h2 className="text-2xl">Create a new expense</h2>
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
                value={expenseName}
                onChange={e => setExpenseName(e.target.value)}
              />
            </div>
            <div className="mt-6">
              <label className="block font-semibold" htmlFor="userID">
                User ID
              </label>
              <input
                className="shadow appearance-none w-full border rounded mt-2 py-2 px-3 leading-tight focus:outline-none focus:ring"
                id="userID"
                name="user_id"
                type="text"
                value={expenseUserID}
                onChange={e => setExpenseUserID(e.target.value)}
              />
            </div>
            <div className="mt-6 flex justify-end">
              <button
                className="bg-indigo-500 hover:bg-indigo-700 text-white py-2 px-4 rounded focus:outline-none focus:ring"
                type="submit"
              >
                Create expense
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
