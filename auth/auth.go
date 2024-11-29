package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauth2Config     *oauth2.Config
	oauthStateString string
	oauth2Token      *oauth2.Token
)

var store = sessions.NewCookieStore([]byte("your-secret-key")) // Replace with a secure secret key

// Initialize OAuth2 configuration with Google
func InitOAuth(clientID, clientSecret, redirectURL string) {
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random" // You can generate a random string here (for better security)
}

// Handle the login request, redirecting the user to Google OAuth
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// Handle the callback from Google, exchanging the code for a token
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != oauthStateString {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to get token from Google: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the token in the session
	oauth2Token = token

	// Get user info from Google API
	client := oauth2Config.Client(r.Context(), oauth2Token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Here, you can extract user information from the response
	// For example: user's email, profile picture, name, etc.

	// Redirect the user to their dashboard after login
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
