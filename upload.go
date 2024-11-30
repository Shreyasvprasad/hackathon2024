package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// FileMetadata structure to hold file information
type FileMetadata struct {
	Filename   string
	FileURL    string
	FileSize   int64
	UploadTime time.Time
}

// In-memory storage for uploaded files
var fileHistory []FileMetadata

// MinIO client
var minioClient *minio.Client

// MinIO configuration
var (
	minioHost      = "localhost"
	minioPort      = 9000
	minioAccessKey = "minioadmin"
	minioSecretKey = "minioadmin"
	bucketName     = "uploads"
)
var oauth2Config = oauth2.Config{
	ClientID:     "1053445861061-ao51cpn5qnu3ajav131jlqqcfsb2bt6s.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8081",
	Scopes:       []string{"openid", "profile", "email"},
	Endpoint:     google.Endpoint,
}

var oauthStateString = "randomstate"

// Initialize MinIO client
func initMinioClient() {
	var err error
	endpoint := fmt.Sprintf("%s:%d", minioHost, minioPort)
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Error initializing MinIO client: %v", err)
	}

	// Check if bucket exists, if not create it
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Error checking if bucket exists: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Error creating bucket: %v", err)
		}
	}
	fmt.Println("MinIO client initialized successfully.")
}

func uploadFileToMinio(file multipart.File, filename string) (string, int64, error) {
	// Create a buffer to determine file size
	fileSizeBuffer := &bytes.Buffer{}

	// Copy file contents to buffer to calculate size
	fileSize, err := io.Copy(fileSizeBuffer, file)
	if err != nil {
		return "", 0, fmt.Errorf("failed to determine file size: %v", err)
	}

	// Reset file pointer to beginning for uploading
	file.Seek(0, io.SeekStart)

	// Upload file to MinIO
	_, err = minioClient.PutObject(context.Background(), bucketName, filename, file, fileSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	// Generate file URL
	fileURL := fmt.Sprintf("http://%s:%d/%s/%s", minioHost, minioPort, bucketName, filename)

	// Store file metadata in memory
	fileHistory = append(fileHistory, FileMetadata{
		Filename:   filename,
		FileURL:    fileURL,
		FileSize:   fileSize,
		UploadTime: time.Now(),
	})

	return fileURL, fileSize, nil
}

// Get file size
func getFileSize(file multipart.File) (int64, error) {
	tempFile, ok := file.(*os.File)
	if !ok {
		return 0, fmt.Errorf("file type is not valid for size extraction")
	}
	fileInfo, err := tempFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("unable to retrieve file stats: %v", err)
	}
	return fileInfo.Size(), nil
}

// Handle file upload form
// Handle file upload form
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := `
		<!DOCTYPE html>
		<html>
		<body>
			<h1>Upload File</h1>
			<form enctype="multipart/form-data" action="/" method="post">
				<input type="file" name="file">
				<button type="submit">Upload</button>
			</form>
			<br>
			<a href="/history">View Upload History</a>
			<br><br>
			<!-- Login to Google Button -->
			<a href="/login"><button>Login to Google</button></a>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, tmpl)
		return
	}

	if r.Method == http.MethodPost {
		// Parse the multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the file from the form
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload file to MinIO
		fileURL, _, err := uploadFileToMinio(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("File upload failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Success response
		fmt.Fprintf(w, `<h1>File uploaded successfully!</h1><p><a href="%s">%s</a></p>`, fileURL, fileURL)
		fmt.Fprint(w, `<button onclick="window.location.href='/'">Return to Upload Page</button>`)
	}
}

// Handle file history request (list all uploaded files)
func historyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Start the HTML output
	fmt.Fprintf(w, "<h1>Uploaded Files</h1>")
	fmt.Fprintf(w, "<table border='1'><tr><th>Filename</th><th>Upload Time</th><th>Size (bytes)</th><th>Download Link</th></tr>")

	// Loop through fileHistory and display metadata
	for _, file := range fileHistory {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%d</td><td><a href='%s'>Download</a></td></tr>",
			file.Filename, file.UploadTime.Format(time.RFC1123), file.FileSize, file.FileURL)
	}

	fmt.Fprintf(w, "</table>")
	fmt.Fprint(w, `<button onclick="window.location.href='/'">Return to Upload Page</button>`)
}

// Start Google OAuth login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// Google OAuth callback handler
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %s", err), http.StatusInternalServerError)
		return
	}

	// Get user information using the token
	client := oauth2Config.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %s", err), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	// You can store user data or set session here
	fmt.Fprintf(w, "User Info: %s", userInfoResp.Body)
}

func main() {
	// Initialize MinIO client
	initMinioClient()

	http.HandleFunc("/", uploadHandler)
	http.HandleFunc("/history", historyHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth/callback", callbackHandler)

	// Start the server
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
