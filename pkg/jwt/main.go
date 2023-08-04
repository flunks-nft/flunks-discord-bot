package jwt

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/flunks-nft/discord-bot/pkg/utils"
)

var (
	JWT_SECRET string
)

func init() {
	utils.LoadEnv()

	JWT_SECRET = os.Getenv("JWT_SECRET")
}

func IsValidJWT(tokenString string) (bool, error) {
	// Split the token string into two parts: header and payload
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("Invalid token format")
	}

	// Parse the token and verify the signature with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	// Check if the parsing was successful
	if err != nil {
		return false, fmt.Errorf("Failed to parse token: %v", err)
	}

	// Check if the token is valid and has not expired
	if token.Valid {
		return true, nil
	}

	return false, nil
}

// RetrieveWalletAddress both parses and verifies the JWT token
// and then extracts the "addr" claim (wallet address) from the token's payload.
// If the token is not valid or does not contain the expected claim,
// appropriate error messages are returned.
func RetrieveWalletAddress(tokenString string) (string, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Replace "YOUR_SECRET_KEY" with the actual secret key used for signing the tokens
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return "", fmt.Errorf("Failed to parse the JWT token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("Invalid JWT token")
	}

	// Access the "addr" claim from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	addr, ok := claims["addr"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid wallet address in the token")
	}

	return addr, nil
}

type jwtKey struct{}
