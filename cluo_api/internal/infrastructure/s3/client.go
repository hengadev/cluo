package s3Storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ClientConfig holds configuration for creating an S3 client
type ClientConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // Optional: for LocalStack or MinIO
}

// NewS3Client creates a new AWS S3 client with the provided configuration
func NewS3Client(ctx context.Context, cfg ClientConfig) (*s3.Client, error) {
	// Load AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"", // Session token (empty for static credentials)
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create S3 client
	var client *s3.Client
	if cfg.Endpoint != "" {
		// Use custom endpoint (for LocalStack, MinIO, etc.)
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Required for LocalStack/MinIO
		})
	} else {
		// Use standard AWS S3
		client = s3.NewFromConfig(awsCfg)
	}

	return client, nil
}

// NewS3ClientFromEnv creates an S3 client using AWS credentials from environment variables
// Expects AWS_REGION, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY to be set
func NewS3ClientFromEnv(ctx context.Context) (*s3.Client, error) {
	// Load default config (uses environment variables)
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg), nil
}
