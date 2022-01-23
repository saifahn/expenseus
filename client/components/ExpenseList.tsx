import { ExpenseAPI } from "api";
import { useState, useRef, useEffect, FormEvent } from "react";
import ExpenseCard from "./ExpenseCard";

export interface Expense {
  name: string;
  id: string;
  userId: string;
  imageUrl?: string;
}

export default function ExpenseList() {
  const [expenses, setExpenses] = useState<Expense[]>();
  const [expenseName, setExpenseName] = useState("");
  const [{ status, error }, setStatus] = useState({
    status: "idle",
    error: null,
  });
  const [statusMessage, setStatusMessage] = useState<string>();
  const cancelled = useRef(false);
  const imageInput = useRef(null);

  async function fetchExpenses() {
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

  async function createExpense(data: FormData) {
    try {
      const api = new ExpenseAPI();
      const response = await api.createExpense(data);
      setStatusMessage(response);
    } catch (err) {
      console.error(err);
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setStatus({ status: "loading", error: null });
    try {
      const data = new FormData();
      data.append("expenseName", expenseName);
      if (imageInput.current?.files.length) {
        data.append("image", imageInput.current.files[0]);
      }

      await createExpense(data);
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
    <section className="p-6 shadow-lg bg-indigo-50 rounded-xl">
      <h2 className="text-2xl">Expenses</h2>
      {expenses &&
        expenses.map(expense => (
          <ExpenseCard expense={expense} key={expense.id} />
        ))}
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
                type="text"
                value={expenseName}
                onChange={e => setExpenseName(e.target.value)}
                required
              />
            </div>
            <div className="mt-6">
              <label className="block font-semibold" htmlFor="addPicture">
                Add a picture?
              </label>
              <input
                id="addPicture"
                type="file"
                role="button"
                aria-label="Add picture"
                accept="image/*"
                ref={imageInput}
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
