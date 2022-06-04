import { useUserContext } from 'context/user';
import Head from 'next/head';
import Link from 'next/link';

export default function Layout({ children }) {
  const { user, error } = useUserContext();
  return (
    <>
      <Head>
        <title>Expenseus</title>
        <meta name="description" content="Generated by create next app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      {!error && !user ? (
        <div>Loading...</div>
      ) : (
        <nav className="container mx-auto mt-4">
          <ul className="flex items-center">
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
            {!error && user && (
              <li className="ml-4">
                <a
                  href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/logout`}
                  className="inline-flex items-center border rounded-md px-3 py-2"
                >
                  <span className="">Log out</span>
                </a>
              </li>
            )}
          </ul>
        </nav>
      )}

      {error && error.code === 401 && (
        <main className="container h-full mx-auto border-gray-600 mt-4">
          <a
            href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/login_google`}
            className="inline-flex items-center border rounded-md px-3 py-2 mt-4"
          >
            <img
              src="/images/google-g-logo.svg"
              alt="Google G Logo"
              height={24}
              width={24}
            />
            <span className="ml-3">Sign in with Google</span>
          </a>
        </main>
      )}

      {error && error.code !== 401 && (
        <>
          <p>There was an error. Please refresh and try again.</p>
        </>
      )}

      {!error && user && (
        <main className="container h-full mx-auto border-4 border-gray-600 mt-4 p-8">
          {children}
        </main>
      )}
    </>
  );
}
