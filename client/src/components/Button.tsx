// const handleClick = () => {
//   // Here you will redirect to the Discord OAuth2 link
//   window.location.href = `https://discord.com/api/oauth2/authorize?client_id=${
//     process.env.REACT_APP_DISCORD_CLIENT_ID
//   }&redirect_uri=${encodeURIComponent(
//     process.env.REACT_APP_REDIRECT_URI
//   )}&response_type=code&scope=identify%20guilds.join`;
// };

const discordUrl =
  "https://discord.com/api/oauth2/authorize?client_id=1121560033600208936&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flogin&response_type=code&scope=identify";

const handleClick = () => {
  // Here you will redirect to the Discord OAuth2 link
  window.location.href = discordUrl;
};

const Button: React.FC<{ text: string }> = ({ text }) => {
  return (
    <div className="mb-7 flex flex-col items-center text-center">
      <div className="flex justify-between gap-2">
        <button
          className="border-orange bg-orange hover:bg-orange-dark rounded-md border px-4 py-2 text-black"
          onClick={handleClick}
        >
          {text}
        </button>
      </div>
    </div>
  );
};

export default Button;
