// @ts-ignore
import jwt from "jsonwebtoken";

// Function to generate and sign the JWT
const generateJWT = (walletAddress: string) => {
  const secretKey = process.env.NEXT_PUBLIC_JWT_SECRET; // Replace with your secret key
  const payload = { addr: walletAddress };
  const options = { expiresIn: 60 * 60 };

  return jwt.sign(payload, secretKey, options);
};

export default generateJWT;
