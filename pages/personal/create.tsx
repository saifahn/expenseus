import PersonalLayout from 'components/LayoutPersonal';
import TxnCreateForm from 'components/TxnCreateForm';
import Head from 'next/head';

export default function PersonalCreate() {
  return (
    <PersonalLayout>
      <Head>
        <title>create personal transaction - expenseus</title>
      </Head>
      <TxnCreateForm />
    </PersonalLayout>
  );
}
