const handleClick = () => {};

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
