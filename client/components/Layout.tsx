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
            <Link href="/personal">
              <a>Personal</a>
            </Link>
          </li>
          <li className="ml-4">
            <Link href="/shared">
              <a>Shared</a>
            </Link>
          </li>
          <li className="ml-4">
            <Link href="/bts">
              <a>BTS</a>
            </Link>
          </li>
        </ul>
      </nav>
      <main className="container h-full mx-auto border-4 border-gray-600 mt-4 p-8">
        {children}
      </main>
    </>
  );
}
