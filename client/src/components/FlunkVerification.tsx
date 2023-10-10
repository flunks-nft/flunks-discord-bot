import React from "react";
import Button from "~/components/Button";
import { useWeb3Context } from "../contexts/Web3";

const FlunkVerification = () => {
  const { user } = useWeb3Context();

  return (
    <div className="w-96 rounded border border-black p-6 text-center">
      <h2 className="mb-4 text-2xl">VERIFY YOUR FLUNKS!</h2>
      <p className="mb-4">
        {user?.loggedIn
          ? `Step 2/2: You'll need to provide Ace authorization to access your Discord Account.`
          : `Step 1/2: We'll need to verify all Flunks in your wallet before they can join the School Yard Battles on Discord.`}
      </p>
      <Button />
    </div>
  );
};

export default FlunkVerification;
