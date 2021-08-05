async function createUser() {
  // make a post request to the URL
}

async function listUsers() {
  // make a request to the URL
  // return users
}

async function handleSubmit() {
  await createUser();
  await listUsers();
}

/**
 * Component rendering a list of users
 */
export default function Users() {
  return (
    <section>
      <h2>Users</h2>
      <form onSubmit={handleSubmit}></form>
    </section>
  );
}
