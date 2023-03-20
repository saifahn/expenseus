import UserList from "components/UserList";

export default function BTS() {
  return (
    <>
      <h1 className="text-4xl">Behind the scenes</h1>
      <p className="mt-4">This is the behind the scenes page</p>
      <div className="mt-4">
        <UserList />
      </div>
    </>
  );
}
