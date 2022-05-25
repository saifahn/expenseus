import 'tailwindcss/tailwind.css';
import Layout from 'components/Layout';
import { UserProvider } from 'context/user';

function MyApp({ Component, pageProps }) {
  return (
    <UserProvider>
      <Layout>
        <Component {...pageProps} />
      </Layout>
    </UserProvider>
  );
}

export default MyApp;
