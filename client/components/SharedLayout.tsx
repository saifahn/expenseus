import Link from 'next/link';

export default function SharedLayout({ children }) {
  return (
    <>
      <h1 className="text-4xl">Shared</h1>
      <nav className="mt-4">
        <ul className="flex">
          <li className="p-2 border-2 cursor-pointer">
            <Link href="/shared">Home</Link>
          </li>
          <li className="p-2 border-2 cursor-pointer ml-4">
            <Link href="/shared/trackers">Trackers</Link>
          </li>
        </ul>
      </nav>
      <section className="mt-4">{children}</section>
    </>
  );
}
