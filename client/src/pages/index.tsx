import Head from "next/head";
import Button from "~/components/Button";

export default function Home() {
  return (
    <>
      <Head>
        <title>Flunks Discord Dapper Portal</title>
        <meta name="description" content="Flunks Discord Dapper Portal" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className="flex min-h-screen flex-col items-center justify-center bg-white">
        <div className="container flex flex-col items-center justify-center gap-12 px-4 py-16 ">
          <h1 className="text-orange text-5xl font-extrabold tracking-tight sm:text-[5rem]">
            Flunks <span className="text-orange">Verification</span> App
          </h1>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-8">
            <Button />
          </div>
        </div>
      </main>
    </>
  );
}
