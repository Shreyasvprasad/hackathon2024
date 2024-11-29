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

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func InitOAuth(clientID, clientSecret, redirectURL string) {
	oauth2Config = &oauth2.Config{
		ClientID:     1053445861061 - ao51cpn5qnu3ajav131jlqqcfsb2bt6s.apps.googleusercontent.com,
		ClientSecret: GOCSPX - S1rpZymdSMI4Tl - x2nOnRL3TL8E5,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random" // You can generate a random string here
}

// OAuth login handler that redirects to Google
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// OAuth callback handler that exchanges the code for a token
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

	// Save the token in the session (store in database or session)
	oauth2Token = token

	// Get user info from Google
	client := oauth2Config.Client(r.Context(), oauth2Token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the user info response here
	// You can extract details like email, name, profile picture, etc.
	// For now, we'll just redirect to the dashboard after successful login
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
