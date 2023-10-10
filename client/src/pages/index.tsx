import Head from "next/head";
import FlunkVerification from "~/components/FlunkVerification";

export default function Home() {
  return (
    <>
      <Head>
        <title>Verify Your Flunks</title>
        <meta name="description" content="Verify Your Flunks" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className="flex min-h-screen flex-col items-center justify-center bg-white">
        <FlunkVerification />
      </main>
    </>
  );
}
