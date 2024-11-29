package auth

import (
	"cloud-storage-app/models"
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauth2Config *oauth2.Config

// InitOAuth initializes OAuth2 with Google credentials
func InitOAuth(clientID, clientSecret, redirectURL string) {
	oauth2Config = &oauth2.Config{
		ClientID:     1053445861061 - ao51cpn5qnu3ajav131jlqqcfsb2bt6s.apps.googleusercontent.com,
		ClientSecret: GOCSPX - S1rpZymdSMI4Tl - x2nOnRL3TL8E5,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleLogin initiates the OAuth flow by redirecting the user to Google
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GoogleCallback handles the callback after OAuth flow
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tok, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	client := oauth2Config.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user models.User
	// Decode response and map to user object (use appropriate JSON parsing here)
	// Assume you get user info with email, name, etc.

	// Save the user to the database (ScyllaDB)
}
