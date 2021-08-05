import Head from "next/head";
import "tailwindcss/tailwind.css";

export default function Home() {
  return (
    <div className="container mx-auto px-4">
      <Head>
        <title>Expenseus</title>
        <meta name="description" content="Generated by create next app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="py-4">
        <h1 className="text-4xl">Welcome to Expenseus</h1>
      </main>
    </div>
  );
}
