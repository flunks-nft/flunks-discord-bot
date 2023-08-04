import { useWeb3Context } from "../contexts/Web3";

const Button: React.FC<{}> = ({}) => {
  const { connect, user, executeScript, logout } = useWeb3Context();

  const handleClick = () => {
    if (!user?.loggedIn) {
      connect();
    } else {
      logout();
    }
  };

  const redirectToDiscord = () => {
    const loginUrl = `http://localhost:8080/auth/login?addr=${user.addr}`;
    window.location.href = loginUrl;
  };

  return (
    <div className="mb-7 flex flex-col items-center text-center">
      <div className="flex justify-between gap-2">
        <button
          className="border-orange bg-orange hover:bg-orange-dark rounded-md border px-4 py-2 text-black"
          onClick={user?.loggedIn ? redirectToDiscord : handleClick}
        >
          {user?.loggedIn ? "Logout" : "Connect Dapper"}
        </button>
      </div>
      <div>
        <h1>{JSON.stringify(user, null, 2)}</h1>
      </div>
    </div>
  );
};

export default Button;
