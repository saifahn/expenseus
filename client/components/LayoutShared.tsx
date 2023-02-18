import { PropsWithChildren } from 'react';

export default function SharedLayout({ children }: PropsWithChildren<{}>) {
  return (
    <>
      {/* <nav className="mt-4">
        <ul className="flex">
          <li className="flex">
            <Link href="/shared">
              <a className="border-2 p-2">Home</a>
            </Link>
          </li>
          <li className="flex">
            <Link href="/shared/trackers">
              <a className="ml-4 border-2 p-2">Trackers</a>
            </Link>
          </li>
        </ul>
      </nav> */}
      <section className="mt-4">{children}</section>
    </>
  );
}
