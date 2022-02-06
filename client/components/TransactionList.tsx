import { TransactionAPI } from "api";
import { useState, useRef, useEffect, FormEvent } from "react";
import TransactionCard from "./TransactionCard";

export interface Transaction {
  name: string;
  id: string;
  userId: string;
  amount: number;
  imageUrl?: string;
}

export default function TransactionList() {
  const [transactions, setTransactions] = useState<Transaction[]>();
  const [transactionName, setTransactionName] = useState("");
  const [amount, setAmount] = useState("");
  const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
  const [{ status, error }, setStatus] = useState({
    status: "idle",
    error: null,
  });
  const [statusMessage, setStatusMessage] = useState<string>();
  const cancelled = useRef(false);
  const imageInput = useRef(null);

  async function fetchTransactions() {
    try {
      const api = new TransactionAPI();
      const transactions = await api.listTransactions();
      if (!cancelled.current) {
        setTransactions(transactions);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function createTransaction(data: FormData) {
    try {
      const api = new TransactionAPI();
      const response = await api.createTransaction(data);
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
      data.append("transactionName", transactionName);
      data.append("amount", amount);
      data.append("date", Date.parse(date).toString());
      if (imageInput.current?.files.length) {
        data.append("image", imageInput.current.files[0]);
      }

      await createTransaction(data);
      setStatus({ status: "fulfilled", error: null });
      await fetchTransactions();
    } catch (err) {
      setStatus({ status: "rejected", error: err });
    }
  }

  useEffect(() => {
    fetchTransactions();
    return () => {
      cancelled.current = true;
    };
  }, []);

  return (
    <section className="p-6 shadow-lg bg-indigo-50 rounded-xl">
      <h2 className="text-2xl">Transactions</h2>
      {transactions &&
        transactions.map(transaction => (
          <TransactionCard transaction={transaction} key={transaction.id} />
        ))}
      <div className="mt-6">
        <h2 className="text-2xl">Create a new transaction</h2>
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
                value={transactionName}
                onChange={e => setTransactionName(e.target.value)}
                required
              />
            </div>
            <div>
              <label className="block font-semibold mt-6" htmlFor="amount">
                Amount
              </label>
              <input
                className="shadow appearance-none w-full border rounded mt-2 py-2 px-3 leading-tight focus:outline-none focus:ring"
                id="amount"
                type="number"
                value={amount}
                onChange={e => setAmount(e.target.value)}
                required
              />
            </div>
            <div>
              <label className="block font-semibold mt-6" htmlFor="date">
                Date
              </label>
              <input
                className="shadow appearance-none w-full border rounded mt-2 py-2 px-3 leading-tight focus:outline-none focus:ring"
                id="amount"
                type="date"
                value={date}
                onChange={e => setDate(e.target.value)}
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
                Create transaction
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
