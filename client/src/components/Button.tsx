import { useWeb3Context } from "../contexts/Web3";
import generateJWT from "./jwt";

const { NEXT_PUBLIC_AUTH_SERVER_URL } = process.env;

const Button: React.FC<Record<string, never>> = () => {
  const { connect, logout, user } = useWeb3Context();

  const handleClick = () => {
    if (!user?.loggedIn) {
      connect();
    } else {
      logout();
    }
  };

  const redirectToDiscord = () => {
    const token = generateJWT(user.addr);
    const loginUrl = `${NEXT_PUBLIC_AUTH_SERVER_URL}auth/login?token=${token}`;
    window.location.href = loginUrl;
  };

  return (
    <div className="mb-7 flex flex-col items-center text-center">
      <div className="flex justify-between gap-2">
        <button
          className="border-orange bg-orange hover:bg-orange-dark rounded-md border px-4 py-2 text-black"
          onClick={user?.loggedIn ? redirectToDiscord : handleClick}
        >
          {user?.loggedIn ? "Click to Verify" : "Connect Dapper"}
        </button>
      </div>
      <div>
        <h1>{JSON.stringify(user, null, 2)}</h1>
      </div>
    </div>
  );
};

export default Button;
