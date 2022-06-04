import 'tailwindcss/tailwind.css';
import Layout from 'components/Layout';
import { UserProvider } from 'context/user';
import { SWRConfig } from 'swr';
import { fetcher } from 'config/fetcher';

function MyApp({ Component, pageProps }) {
  return (
    <SWRConfig value={{ fetcher }}>
      <UserProvider>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </UserProvider>
    </SWRConfig>
  );
}

export default MyApp;
