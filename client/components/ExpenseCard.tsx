import Image from "next/image";
import { Expense } from "./ExpenseList";

export default function ExpenseCard({ expense }: { expense: Expense }) {
  return (
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
          width={400}
          height={400}
          objectFit="contain"
          alt="expense image"
        />
      )}
    </article>
  );
}
