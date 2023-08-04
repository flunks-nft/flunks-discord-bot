import jwt from "jsonwebtoken";

// Function to generate and sign the JWT
const generateJWT = (walletAddress: string) => {
  console.log("--NEXT_PUBLIC_JWT_SECRET", process.env.NEXT_PUBLIC_JWT_SECRET);

  const secretKey = process.env.NEXT_PUBLIC_JWT_SECRET; // Replace with your secret key
  const payload = { addr: walletAddress };
  const options = { expiresIn: 60 * 60 };

  console.log("--payload", payload);
  console.log("--secretKey", secretKey);
  console.log("--options", options);

  return jwt.sign(payload, secretKey, options);
};

export default generateJWT;
