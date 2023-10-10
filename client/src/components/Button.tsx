import { useWeb3Context } from "../contexts/Web3";
import generateJWT from "./jwt";

const Button: React.FC<Record<string, never>> = () => {
  const { connect, logout, user } = useWeb3Context();

  const handleLoginToggle = () => {
    if (!user?.loggedIn) {
      connect();
    } else {
      logout();
    }
  };

  const redirectToDiscord = () => {
    const token = generateJWT(user.addr);
    const loginUrl = `https://discord-runner-s2ncmw3esa-uw.a.run.app/auth/login?token=${token}`;
    window.location.href = loginUrl;
  };

  return (
    <div className="mb-7 flex flex-col items-center text-center">
      <div className="flex flex-col justify-between gap-2">
        <button
          className="rounded bg-black px-4 py-2 text-white transition duration-300 hover:bg-gray-600"
          onClick={user?.loggedIn ? redirectToDiscord : handleLoginToggle}
        >
          {user?.loggedIn ? "GOT IT" : "Connect Wallet"}
        </button>
        {user?.loggedIn && (
          <button
            className="rounded bg-black px-4 py-2 text-white transition duration-300 hover:bg-gray-600"
            onClick={user?.loggedIn ? handleLoginToggle : () => {}}
          >
            {"Sign Out"}
          </button>
        )}
      </div>
    </div>
  );
};

export default Button;
