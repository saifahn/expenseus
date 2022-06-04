import Link from 'next/link';

export default function SharedLayout({ children }) {
  return (
    <>
      <h1 className="text-4xl">Shared</h1>
      <nav className="mt-4">
        <ul className="flex">
          <li>
            <Link href="/shared">
              <a className="p-2 border-2">Home</a>
            </Link>
          </li>
          <li>
            <Link href="/shared/trackers">
              <a className="p-2 border-2 ml-4">Trackers</a>
            </Link>
          </li>
        </ul>
      </nav>
      <section className="mt-4">{children}</section>
    </>
  );
}
