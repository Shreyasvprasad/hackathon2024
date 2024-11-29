package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callback",
	ClientID:     "1053445861061-ao51cpn5qnu3ajav131jlqqcfsb2bt6s.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var notes []Note

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	fmt.Println("File name:", header.Filename)
	defer file.Close()

	tempFile, err := ioutil.TempFile("uploads", "upload-*.png")
	if err != nil {
		fmt.Println("Error creating the file:", err)
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully uploaded file")
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Error decoding the note", http.StatusBadRequest)
		return
	}

	notes = append(notes, note)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func main() {
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/note", createNoteHandler)
	http.ListenAndServe(":8080", nil)
}
