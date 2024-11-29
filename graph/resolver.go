package resolvers

import (
	"cloud-storage-app/models"
	"context"
)

func (r *mutationResolver) CreateNote(ctx context.Context, userID string, title string, content string) (*models.Note, error) {
	// Insert the note into ScyllaDB
	// Return the created note
}

func (r *mutationResolver) UploadFile(ctx context.Context, userID string, file models.File) (*models.File, error) {
	// Upload the file to MinIO or S3 and store file metadata in ScyllaDB
	// Return the uploaded file metadata
}

func (r *queryResolver) GetNotes(ctx context.Context, userID string) ([]*models.Note, error) {
	// Retrieve the user's notes from ScyllaDB
}
