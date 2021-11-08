import { ExpenseAPI } from "api";
import { Expense } from "components/ExpenseList";
import { useRouter } from "next/router";
import { useEffect, useState, useRef } from "react";
import Image from "next/image";

const SingleExpense = () => {
  const router = useRouter();
  const { eid } = router.query;
  const [expense, setExpense] = useState<Expense>();
  const cancelled = useRef(false);

  async function fetchExpense(id: string) {
    try {
      const api = new ExpenseAPI();
      const expense = await api.getExpense(id);
      if (!cancelled.current) {
        setExpense(expense);
      }
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (!router.isReady) return;
    fetchExpense(eid as string);
    return () => {
      cancelled.current = true;
    };
  });

  return (
    <main className="container">
      <h2>Expense name: {eid}</h2>
      {expense && (
        <article
          className="p-4 mt-4 rounded-md shadow-md bg-white"
          key={expense.id}
        >
          <h3 className="text-xl">{expense.name}</h3>
          <p>{expense.userID}</p>
          <p>{expense.id}</p>
          {expense.imageURL && (
            <Image
              src={expense.imageURL}
              layout="fill"
              objectFit="contain"
              alt="expense image"
            />
          )}
        </article>
      )}
    </main>
  );
};

export default SingleExpense;
