import jwt from "jsonwebtoken";

// Function to generate and sign the JWT
const generateJWT = (walletAddress: string) => {
  const secretKey = process.env.REACT_APP_JWT_SECRET; // Replace with your secret key
  const payload = { addr: walletAddress };
  const options = { expiresIn: "5 minutes" }; // Set an appropriate expiration time

  return jwt.sign(payload, secretKey, options);
};

export default generateJWT;
