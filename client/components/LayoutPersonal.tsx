import Link from 'next/link';

export default function PersonalLayout({ children }) {
  return (
    <>
      <h1 className="text-4xl">Personal</h1>
      <nav className="mt-4">
        <ul className="flex">
          <li className="flex">
            <Link href="/personal">
              <a className="border-2 p-2">Home</a>
            </Link>
          </li>
          <li className="flex">
            <Link href="/personal/analysis">
              <a className="ml-4 border-2 p-2">Analysis</a>
            </Link>
          </li>
        </ul>
        <section className="mt-4">{children}</section>
      </nav>
    </>
  );
}
