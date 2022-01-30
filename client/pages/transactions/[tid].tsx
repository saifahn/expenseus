import { useRouter } from "next/router";
import { useEffect, useState, useRef } from "react";
import { TransactionAPI } from "api";
import { Transaction } from "components/TransactionList";
import TransactionCard from "components/TransactionCard";

const SingleTransaction = () => {
  const router = useRouter();
  const { tid } = router.query;
  const [transaction, setTransaction] = useState<Transaction>();
  const [{ status, error }, setStatus] = useState({
    status: "idle",
    error: null,
  });
  const cancelled = useRef(false);

  async function fetchTransaction(id: string) {
    try {
      setStatus({ status: "pending", error: null });
      const api = new TransactionAPI();
      const transaction = await api.getTransaction(id);
      if (!cancelled.current) {
        setTransaction(transaction);
      }
      setStatus({ status: "fulfilled", error: null });
    } catch (err) {
      setStatus({ status: "rejected", error: err });
    }
  }

  useEffect(() => {
    if (!router.isReady) return;
    fetchTransaction(tid as string);
    return () => {
      cancelled.current = true;
    };
  }, [tid, router.isReady]);

  return (
    <main className="container">
      {error ? (
        <h4>Sorry, there was an error, please try again</h4>
      ) : status === "fulfilled" ? (
        transaction ? (
          <TransactionCard transaction={transaction} />
        ) : (
          <h4>Sorry, no transaction found for ID: {tid}</h4>
        )
      ) : (
        <h4>loading</h4>
      )}
    </main>
  );
};

export default SingleTransaction;
