import Link from 'next/link';
import { PropsWithChildren } from 'react';

export default function PersonalLayout({ children }: PropsWithChildren<{}>) {
  return (
    <>
      <nav className="mt-4">
        <ul className="flex">
          <li className="mr-4 flex flex-1">
            <Link href="/personal/analysis">
              <a className="w-full rounded-lg bg-emerald-50 py-3 px-4 font-medium lowercase text-black hover:bg-emerald-100 active:bg-emerald-200">
                🔎 Analyze
              </a>
            </Link>
          </li>
          <li className="flex flex-1">
            <Link href="/personal/create">
              <a className="w-full rounded-lg bg-violet-50 py-3 px-4 font-medium lowercase text-black hover:bg-violet-100 active:bg-violet-200">
                ➕ Create new
              </a>
            </Link>
          </li>
        </ul>
        <section className="mt-4">{children}</section>
      </nav>
    </>
  );
}
