import { useRouter } from "next/router";
import { useEffect, useState, useRef } from "react";
import { ExpenseAPI } from "api";
import { Expense } from "components/ExpenseList";
import ExpenseCard from "components/ExpenseCard";

const SingleExpense = () => {
  const router = useRouter();
  const { eid } = router.query;
  const [expense, setExpense] = useState<Expense>();
  const [{ status, error }, setStatus] = useState({
    status: "idle",
    error: null,
  });
  const cancelled = useRef(false);

  async function fetchExpense(id: string) {
    try {
      setStatus({ status: "pending", error: null });
      const api = new ExpenseAPI();
      const expense = await api.getExpense(id);
      if (!cancelled.current) {
        setExpense(expense);
      }
      setStatus({ status: "fulfilled", error: null });
    } catch (err) {
      setStatus({ status: "rejected", error: err });
    }
  }

  useEffect(() => {
    if (!router.isReady) return;
    fetchExpense(eid as string);
    return () => {
      cancelled.current = true;
    };
  }, [eid, router.isReady]);

  return (
    <main className="container">
      {error ? (
        <h4>Sorry, there was an error, please try again</h4>
      ) : status === "fulfilled" ? (
        expense ? (
          <ExpenseCard expense={expense} />
        ) : (
          <h4>Sorry, no expense found for ID: {eid}</h4>
        )
      ) : (
        <h4>loading</h4>
      )}
    </main>
  );
};

export default SingleExpense;
