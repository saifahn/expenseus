import Link from "next/link";

export default function Layout({ children }) {
  return (
    <>
      <nav className="container mx-auto mt-4">
        <ul className="flex">
          <li>
            <Link href="/">
              <a>Home</a>
            </Link>
          </li>
          <li className="ml-4">
            <Link href="/users">
              <a>Users</a>
            </Link>
          </li>
          <li className="ml-4">
            <Link href="/expenses">
              <a>Expenses</a>
            </Link>
          </li>
        </ul>
      </nav>
      <main className="container mx-auto mt-4 py-8">{children}</main>
    </>
  );
}
