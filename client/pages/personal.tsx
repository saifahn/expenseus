import TransactionSubmitForm from 'components/TransactionSubmitForm';
import { useUserContext } from 'context/user';
import useSWR from 'swr';

interface Transaction {
  name: string;
  id: string;
  userId: string;
  amount: number;
  imageUrl?: string;
}

export default function Personal() {
  const { user } = useUserContext();
  const { data: transactions, error } = useSWR<Transaction[]>(() =>
    user
      ? `${process.env.NEXT_PUBLIC_API_BASE_URL}/transactions/user/${user.id}`
      : null,
  );

  return (
    <>
      <h1 className="text-4xl">Personal</h1>
      <div className="mt-4">
        <TransactionSubmitForm />
      </div>
      <div className="mt-4 p-4">
        <h2 className="text-2xl">Personal transactions</h2>
        {error && <div>Failed to load transactions</div>}
        {transactions === null && <div>Loading list of transactions...</div>}
        {transactions && transactions.length === 0 && (
          <div>No transactions to show</div>
        )}
        {transactions &&
          transactions.map((txn) => (
            <article className="p-2 border-2 mt-4" key={txn.id}>
              <h3 className="text-lg">{txn.name}</h3>
              <p>{txn.amount}</p>
              <p>{txn.userId}</p>
              <p>{txn.id}</p>
            </article>
          ))}
      </div>
    </>
  );
}
