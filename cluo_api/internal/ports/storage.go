package ports

import (
	"context"
	"io"
)

// StorageService defines operations for file storage (S3, GCS, Azure Blob, etc.)
type StorageService interface {
	// UploadFile uploads a file to storage and returns the public URL
	UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string, fileSize int64) (string, error)

	// DeleteFile removes a file from storage by its URL
	DeleteFile(ctx context.Context, fileURL string) error

	// GetFileURL generates a signed URL for accessing a file (for private buckets)
	// Returns the URL as-is if the bucket is public
	GetFileURL(ctx context.Context, fileURL string) (string, error)
}
