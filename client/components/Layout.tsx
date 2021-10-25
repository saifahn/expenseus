import Link from "next/link";

export default function Layout({ children }) {
  return (
    <>
      <nav>
        <ul>
          <li>
            <Link href="/">
              <a>Home</a>
            </Link>
          </li>
          <li>
            <Link href="/users">
              <a>Users</a>
            </Link>
          </li>
          <li>
            <Link href="/expenses">
              <a>Expenses</a>
            </Link>
          </li>
        </ul>
      </nav>
      <main>{children}</main>
    </>
  );
}
