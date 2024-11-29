package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

// Global ScyllaDB session
var session *gocql.Session

// Initialize ScyllaDB connection
func InitScyllaDB() {
	cluster := gocql.NewCluster("localhost:9042") // Replace with your ScyllaDB host
	cluster.Keyspace = "some_scylla"              // Replace with your keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second

	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	log.Println("Connected to ScyllaDB")
}

// Close ScyllaDB connection
func CloseScyllaDB() {
	if session != nil {
		session.Close()
		log.Println("ScyllaDB connection closed")
	}
}

// Create a new note
func CreateNoteHandler(c *gin.Context) {
	var note struct {
		UserID  gocql.UUID `json:"user_id"`
		Title   string     `json:"title"`
		Content string     `json:"content"`
	}

	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	noteID := gocql.TimeUUID()
	now := time.Now()

	err := session.Query(`INSERT INTO notes (note_id, user_id, title, content, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`, noteID, note.UserID, note.Title, note.Content, now, now).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note created", "note_id": noteID})
}

// Retrieve all notes for a user
func GetNotesHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var notes []map[string]interface{}
	iter := session.Query(`SELECT note_id, title, content, created_at, updated_at FROM notes WHERE user_id = ?`, userID).Iter()

	var noteID gocql.UUID
	var title, content string
	var createdAt, updatedAt time.Time
	for iter.Scan(&noteID, &title, &content, &createdAt, &updatedAt) {
		notes = append(notes, map[string]interface{}{
			"note_id":    noteID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// Update an existing note
func UpdateNoteHandler(c *gin.Context) {
	var note struct {
		NoteID  gocql.UUID `json:"note_id"`
		Title   string     `json:"title"`
		Content string     `json:"content"`
	}

	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	now := time.Now()
	err := session.Query(`UPDATE notes SET title = ?, content = ?, updated_at = ? WHERE note_id = ?`,
		note.Title, note.Content, now, note.NoteID).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated"})
}

// Delete a note
func DeleteNoteHandler(c *gin.Context) {
	noteID := c.Param("note_id")
	if noteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note ID is required"})
		return
	}

	err := session.Query(`DELETE FROM notes WHERE note_id = ?`, noteID).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

func main() {
	// Initialize ScyllaDB connection
	InitScyllaDB()
	defer CloseScyllaDB()

	// Set up Gin router
	router := gin.Default()

	// CRUD routes for notes
	router.POST("/notes", CreateNoteHandler)
	router.GET("/notes", GetNotesHandler)
	router.PUT("/notes", UpdateNoteHandler)
	router.DELETE("/notes/:note_id", DeleteNoteHandler)

	// Start the server
	router.Run(":8080")
}
