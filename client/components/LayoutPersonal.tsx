import Link from 'next/link';

export default function PersonalLayout({ children }) {
  return (
    <>
      <nav className="mt-4">
        <ul className="flex">
          <li className="flex">
            <Link href="/personal/create">
              <a className="rounded-md bg-violet-100 p-2 lowercase hover:bg-violet-200 active:bg-violet-300">
                Create +
              </a>
            </Link>
          </li>
          <li className="flex">
            <Link href="/personal/analysis">
              <a className="ml-4 rounded-md p-2 lowercase hover:bg-slate-200 active:bg-slate-300">
                Analyze
              </a>
            </Link>
          </li>
        </ul>
        <section className="mt-4">{children}</section>
      </nav>
    </>
  );
}
