import { TransactionAPI } from "api";
import { useState, useRef, useEffect, FormEvent } from "react";
import TransactionCard from "./TransactionCard";
import TransactionSubmitForm from "./TransactionSubmitForm";

export interface Transaction {
  name: string;
  id: string;
  userId: string;
  amount: number;
  imageUrl?: string;
}

export default function TransactionList() {
  const [transactions, setTransactions] = useState<Transaction[]>();
  const cancelled = useRef(false);

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
      <TransactionSubmitForm handleAfterSubmit={fetchTransactions} />
    </section>
  );
}
