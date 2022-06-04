import { useUserContext } from 'context/user';

export default function Home() {
  const { user } = useUserContext();
  return (
    <>
      <div>Hi, {user.username}!</div>
      <div>
        You should be able to see a summary of your recent transactions here -
        up to 3 months or 100 transactions, whichever is less.
      </div>
    </>
  );
}
