import SharedLayout from 'components/LayoutShared';
import TrackersSubmitForm from 'components/TrackersSubmitForm';
import Head from 'next/head';

export default function SharedCreateTracker() {
  return (
    <SharedLayout>
      <Head>
        <title>create tracker - expenseus</title>
      </Head>
      <TrackersSubmitForm />
    </SharedLayout>
  );
}
