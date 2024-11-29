package main

import ( // Adjust the import path based on your module name
	"auth"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/callback", auth.CallbackHandler)

	fmt.Println("Server is running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Login with Google</a></body></html>`
	fmt.Fprint(w, html)
}
