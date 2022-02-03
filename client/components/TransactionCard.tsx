import { Transaction } from "./TransactionList";

export default function TransactionCard({
  transaction,
}: {
  transaction: Transaction;
}) {
  return (
    <article
      className="p-4 mt-4 rounded-md shadow-md bg-white"
      key={transaction.id}
    >
      <h3 className="text-xl">{transaction.name}</h3>
      <p>¥{transaction.amount}</p>
      <p>{transaction.userId}</p>
      {transaction.imageUrl && (
        <img
          src={transaction.imageUrl}
          width={400}
          height={400}
          alt="transaction image"
        />
      )}
    </article>
  );
}
