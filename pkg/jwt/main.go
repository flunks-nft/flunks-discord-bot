package jwt

import (
	"fmt"
	"net/http"
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

// Define your secret key for signing and verifying tokens
var secretKey = []byte(JWT_SECRET) // Replace with your secret key

func verifyJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Authorization" header
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := bearerToken[1]

		// Parse and verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, continue with the next handler
		next(w, r)
	}
}

// // secureHandler gets the user's wallet address from the token claims and return wallet address
// func RetrieveWalletAddress(w http.ResponseWriter, r *http.Request) (string, error) {
// 	// Access the user's wallet address from the token claims
// 	claims, ok := r.Context().Value(jwtKey{}).(jwt.MapClaims)
// 	if !ok {
// 		http.Error(w, "Invalid token claims", http.StatusInternalServerError)
// 		return "", nil
// 	}

// 	addr, ok := claims["addr"].(string)
// 	if !ok {
// 		http.Error(w, "Invalid wallet address", http.StatusInternalServerError)
// 		return "", nil
// 	}

// 	return addr, nil
// }

func RetrieveWalletAddress(tokenString string) (string, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Replace "YOUR_SECRET_KEY" with the actual secret key used for signing the tokens
		return []byte("YOUR_SECRET_KEY"), nil
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
