package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/flunks-nft/discord-bot/pkg/db"
	"github.com/flunks-nft/discord-bot/pkg/jwt"
	"github.com/flunks-nft/discord-bot/pkg/utils"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

const (
	discordAPIURL   = "https://discord.com/api"
	discordAuthURL  = discordAPIURL + "/oauth2/authorize"
	discordTokenURL = discordAPIURL + "/oauth2/token"

	discordScopes = "identify" // You can request additional scopes separated by space if needed
	// Generate and store a random state value to prevent CSRF attacks
	STATE_SEED = "FLUNKS_DUNK_STATE"
)

var (
	discordClientID     string
	discordClientSecret string
	discordRedirectURL  string // Discord callback URL
	discordOauth2Config oauth2.Config
)

func init() {
	utils.LoadEnv()

	discordClientID = os.Getenv("DISCORD_CLIENT_ID")
	discordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
	discordRedirectURL = os.Getenv("DISCORD_REDIRECT_URL")

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
}

func StartOauthServer() {
	r := mux.NewRouter()
	r.HandleFunc("/auth/login", handleLogin)
	r.HandleFunc("/auth/callback", handleCallback)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Cloud Run!"))
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// handleLogin sends a user to the Discord login page
// and redirects the user to /auth/callback when authorized.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Get the token from the query parameters
	tokenString := r.URL.Query().Get("token")

	authURL := discordOauth2Config.AuthCodeURL(
		tokenString,
		oauth2.AccessTypeOnline,
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleCallback handles callbacks from the Discord OAuth2 server
// and exchanges the user's information from Discord server with their access token.
func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No authorization code found", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")

	// Validate state
	ok, err := jwt.IsValidJWT(state)
	if err != nil || !ok {
		log.Println("Error in decoding jwt token", err.Error())
		http.Error(w, "Invalid state token", http.StatusBadRequest)
		return
	}

	walletAddress, err := jwt.RetrieveWalletAddress(state)
	if err != nil {
		log.Println("Error in decoding jwt token", err.Error())
		http.Error(w, "Invalid state token", http.StatusBadRequest)
		return
	}

	token, err := discordOauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Use the access token to make API requests on behalf of the user
	client := discordOauth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://discord.com/api/v8/users/@me")
	if err != nil {
		http.Error(w, "Failed to make API request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Define a struct to hold user data
	type User struct {
		ID string `json:"id"`
		// Other fields can be added here as needed
	}

	// Parse the JSON response into the User struct
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Failed to parse JSON response", http.StatusInternalServerError)
		return
	}

	// TODO: write the user's discord id to the database
	err = db.CreateOrUpdateFlowAddress(user.ID, walletAddress)
	if err != nil {
		log.Println("Error while creating or updating flow address in database", err.Error())
		http.Error(w, "Database operation failed", http.StatusInternalServerError)
		return
	}

	// // Return a success message
	// fmt.Fprintln(w, "Authentication successful! Your Wallet is set to:", walletAddress)

	// Render the HTML template
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	tmplPath := filepath.Join(wd, "templates", "authorized.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
