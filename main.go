package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

const (
	discordAPIURL       = "https://discord.com/api"
	discordAuthURL      = discordAPIURL + "/oauth2/authorize"
	discordTokenURL     = discordAPIURL + "/oauth2/token"
	discordClientID     = "1121560033600208936"
	discordClientSecret = "d67EY4CbaPyfX58pS1OLLMM6Yu6LYp3h"
	discordRedirectURL  = "http://localhost:3000/" // Your callback URL
	discordScopes       = "identify"               // You can request additional scopes separated by space if needed
)

var (
	discordOauth2Config = oauth2.Config{
		ClientID:     discordClientID,
		ClientSecret: discordClientSecret,
		RedirectURL:  discordRedirectURL,
		Scopes:       []string{discordScopes},
		Endpoint: oauth2.Endpoint{
			AuthURL:  discordAuthURL,
			TokenURL: discordTokenURL,
		},
	}
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", handleLogin)
	r.HandleFunc("/callback", handleCallback)
	http.Handle("/", r)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate and store a random state value to prevent CSRF attacks
	state := "random_state_value"
	authURL := discordOauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No authorization code found", http.StatusBadRequest)
		return
	}

	// Ensure the state value matches to prevent CSRF attacks
	state := r.URL.Query().Get("state")
	if state != "random_state_value" {
		http.Error(w, "Invalid state value", http.StatusBadRequest)
		return
	}

	token, err := discordOauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Use the access token to make API requests on behalf of the user
	// You can also fetch user details like username, discriminator, avatar, etc.
	// using the access token from the Discord API.

	// For example:
	// client := discordOauth2Config.Client(r.Context(), token)
	// resp, err := client.Get("https://discord.com/api/v10/users/@me")
	// ...

	// You can also store the access token securely and use it to authenticate the user for future requests.
	// Note: Be cautious about how you store and use access tokens as they grant access to the user's account.

	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
}
