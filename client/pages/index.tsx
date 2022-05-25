import { useUserContext } from 'context/user';

export default function Home() {
  const { user } = useUserContext();
  return (
    <>
      <div>Hi, {user.username}!</div>
      <div>Here, we're going to show all of your transactions.</div>
    </>
  );
}
