package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"

	/*"hackathon2024/auth"
	"hackathon2024/db"
	"hackathon2024/graphql"
	"hackathon2024/realtime"
	"hackathon2024/storage"*/
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load(".env")

	InitDB()
	//InitMinIO()

	// Set up routes
	http.HandleFunc("/auth/login", LoginHandler)
	http.HandleFunc("/auth/callback", CallbackHandler)
	http.HandleFunc("/ws/notes", NotesSyncHandler)

	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	//http.Handle("/graphql", handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}})))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// AUTH
var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("randomState", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]string
	json.NewDecoder(resp.Body).Decode(&userInfo)

	fmt.Fprintf(w, "User Info: %+v\n", userInfo)
}

// DB
var Session *gocql.Session

func InitDB() {
	cluster := gocql.NewCluster("172.17.0.2")
	cluster.Keyspace = `scylla`
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to ScyllaDB:", err)
	}
	Session = session
}

//websockets

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NotesSyncHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	for {
		var message string
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		err = conn.WriteJSON(message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
