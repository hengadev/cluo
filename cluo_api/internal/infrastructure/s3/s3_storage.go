package s3Storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/ports"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
	region     string
	baseURL    string // For public buckets, e.g., "https://bucket.s3.region.amazonaws.com"
}

// Config holds the configuration for S3 storage
type Config struct {
	Client     *s3.Client
	BucketName string
	Region     string
	BaseURL    string // Optional: for custom domains or public bucket URLs
}

// New creates a new S3 storage service
func New(cfg Config) ports.StorageService {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		// Default to standard S3 URL format
		baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.BucketName, cfg.Region)
	}

	return &S3Storage{
		client:     cfg.Client,
		bucketName: cfg.BucketName,
		region:     cfg.Region,
		baseURL:    baseURL,
	}
}

// UploadFile uploads a file to S3 and returns the public URL
func (s *S3Storage) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string, fileSize int64) (string, error) {
	// Generate unique key for the file
	key := s.generateFileKey(fileName)

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
		// Optional: Set ACL to public-read if you want files to be publicly accessible
		// ACL: types.ObjectCannedACLPublicRead,
		// Or keep private and use presigned URLs
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the file URL
	fileURL := fmt.Sprintf("%s/%s", s.baseURL, key)
	return fileURL, nil
}

// DeleteFile removes a file from S3 by its URL
func (s *S3Storage) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract key from URL
	key, err := s.extractKeyFromURL(fileURL)
	if err != nil {
		return fmt.Errorf("failed to extract key from URL: %w", err)
	}

	// Delete from S3
	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetFileURL generates a signed URL for accessing a file (for private buckets)
// If the bucket is public, it returns the URL as-is
func (s *S3Storage) GetFileURL(ctx context.Context, fileURL string) (string, error) {
	// For public buckets, return the URL as-is
	// For private buckets, generate a presigned URL

	// Extract key from URL
	key, err := s.extractKeyFromURL(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to extract key from URL: %w", err)
	}

	// Create a presigned URL that expires in 1 hour
	presignClient := s3.NewPresignClient(s.client)
	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Hour * 1 // URL valid for 1 hour
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignResult.URL, nil
}

// generateFileKey generates a unique S3 key for the file
// Format: media/YYYY/MM/DD/{uuid}_{original_filename}
func (s *S3Storage) generateFileKey(fileName string) string {
	now := time.Now()

	// Extract file extension
	ext := filepath.Ext(fileName)

	// Generate unique ID
	uniqueID := uuid.New().String()

	// Sanitize original filename (remove extension, clean special chars)
	baseName := strings.TrimSuffix(fileName, ext)
	baseName = sanitizeFileName(baseName)

	// Create key with date-based folder structure
	key := fmt.Sprintf("media/%d/%02d/%02d/%s_%s%s",
		now.Year(),
		now.Month(),
		now.Day(),
		uniqueID,
		baseName,
		ext,
	)

	return key
}

// extractKeyFromURL extracts the S3 key from a full URL
func (s *S3Storage) extractKeyFromURL(fileURL string) (string, error) {
	// Remove the base URL to get the key
	if !strings.HasPrefix(fileURL, s.baseURL) {
		return "", fmt.Errorf("URL does not match bucket base URL")
	}

	key := strings.TrimPrefix(fileURL, s.baseURL+"/")
	if key == "" {
		return "", fmt.Errorf("invalid file URL: no key found")
	}

	return key, nil
}

// sanitizeFileName removes special characters from filename
func sanitizeFileName(fileName string) string {
	// Replace spaces with underscores
	fileName = strings.ReplaceAll(fileName, " ", "_")

	// Remove or replace other special characters
	replacer := strings.NewReplacer(
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"{", "",
		"}", "",
		"'", "",
		"\"", "",
		",", "",
		";", "",
	)
	fileName = replacer.Replace(fileName)

	// Limit length
	if len(fileName) > 100 {
		fileName = fileName[:100]
	}

	return fileName
}

// DownloadFile downloads a file from S3 by its URL and returns an io.ReadCloser.
// The caller is responsible for closing the reader.
func (s *S3Storage) DownloadFile(ctx context.Context, fileURL string) (io.ReadCloser, error) {
	key, err := s.extractKeyFromURL(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract key from URL: %w", err)
	}

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return resp.Body, nil
}
