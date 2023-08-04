import { type AppType } from "next/dist/shared/lib/utils";
import { Web3ContextProvider } from "~/contexts/Web3";
import "~/styles/globals.css";

const MyApp: AppType = ({ Component, pageProps }) => {
  return (
    <Web3ContextProvider>
      <Component {...pageProps} />
    </Web3ContextProvider>
  );
};

export default MyApp;
