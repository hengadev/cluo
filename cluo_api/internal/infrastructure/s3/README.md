# S3 Storage Service

This package provides S3 storage implementation for media file uploads, deletions, and access.

## Features

- ✅ Upload files to S3 with automatic key generation
- ✅ Date-based folder structure (media/YYYY/MM/DD/)
- ✅ Unique file naming to prevent collisions
- ✅ Delete files from S3
- ✅ Generate presigned URLs for private buckets
- ✅ Support for custom endpoints (LocalStack, MinIO)

## Usage

### 1. Initialize S3 Client

#### Using Environment Variables
```go
import (
    "context"
    s3Storage "github.com/hengadev/cluo_api/internal/infrastructure/s3"
)

// Requires: AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY in environment
client, err := s3Storage.NewS3ClientFromEnv(context.Background())
if err != nil {
    log.Fatal(err)
}
```

#### Using Explicit Configuration
```go
client, err := s3Storage.NewS3Client(context.Background(), s3Storage.ClientConfig{
    Region:          "us-east-1",
    AccessKeyID:     "your-access-key",
    SecretAccessKey: "your-secret-key",
})
```

#### Using LocalStack/MinIO (Development)
```go
client, err := s3Storage.NewS3Client(context.Background(), s3Storage.ClientConfig{
    Region:          "us-east-1",
    AccessKeyID:     "test",
    SecretAccessKey: "test",
    Endpoint:        "http://localhost:4566", // LocalStack
})
```

### 2. Create Storage Service

```go
storage := s3Storage.New(s3Storage.Config{
    Client:     client,
    BucketName: "my-media-bucket",
    Region:     "us-east-1",
    BaseURL:    "", // Optional: leave empty for default S3 URL
})
```

#### With Custom Domain/CDN
```go
storage := s3Storage.New(s3Storage.Config{
    Client:     client,
    BucketName: "my-media-bucket",
    Region:     "us-east-1",
    BaseURL:    "https://cdn.example.com", // Custom domain
})
```

### 3. Use with Media Service

```go
import (
    mediaService "github.com/hengadev/cluo_api/internal/application/media"
)

// Initialize media service with S3 storage
mediaSvc := mediaService.New(
    mediaRepo,
    caseRepo,
    storage,  // S3 storage service
    crypto,
)
```

## File Organization

Files are uploaded with the following structure:
```
media/
  └── 2025/
      └── 01/
          └── 15/
              ├── {uuid}_image1.jpg
              ├── {uuid}_video2.mp4
              └── {uuid}_audio3.mp3
```

Example key: `media/2025/01/15/550e8400-e29b-41d4-a716-446655440000_vacation_photo.jpg`

## S3 Bucket Configuration

### Required Bucket Permissions

Your S3 bucket policy or IAM user must have these permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    }
  ]
}
```

### CORS Configuration (for browser uploads)

```json
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
    "AllowedOrigins": ["https://your-domain.com"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]
```

### Public vs Private Buckets

#### Public Bucket (files accessible via URL)
- Uncomment the `ACL` line in `s3_storage.go`:
  ```go
  ACL: types.ObjectCannedACLPublicRead,
  ```
- Add bucket policy:
  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": "*",
        "Action": "s3:GetObject",
        "Resource": "arn:aws:s3:::your-bucket-name/*"
      }
    ]
  }
  ```

#### Private Bucket (files accessed via presigned URLs)
- Keep bucket private
- Use `GetFileURL()` to generate presigned URLs
- URLs expire after 1 hour (configurable in code)

## Environment Variables

```bash
# AWS Credentials
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key

# S3 Configuration
S3_BUCKET_NAME=my-media-bucket
S3_BASE_URL=https://my-media-bucket.s3.us-east-1.amazonaws.com  # Optional

# For LocalStack (development)
AWS_ENDPOINT=http://localhost:4566
```

## Testing

### Create Test Bucket in LocalStack
```bash
# Start LocalStack
docker run -d -p 4566:4566 localstack/localstack

# Create bucket
aws --endpoint-url=http://localhost:4566 s3 mb s3://test-bucket
```

### Run Tests
```bash
go test ./internal/infrastructure/s3/...
```

## Error Handling

The storage service returns descriptive errors:
- File upload failures
- Invalid URLs
- Key extraction errors
- S3 API errors

Always wrap errors with context:
```go
fileURL, err := storage.UploadFile(ctx, file, fileName, contentType, fileSize)
if err != nil {
    return fmt.Errorf("failed to upload media file: %w", err)
}
```

## Security Considerations

1. **Credentials**: Never commit AWS credentials to version control
2. **Bucket Access**: Use least-privilege IAM policies
3. **Presigned URLs**: Set appropriate expiration times
4. **File Validation**: Validate file types before upload
5. **Rate Limiting**: Implement rate limiting for uploads
6. **Virus Scanning**: Consider integrating virus scanning for uploads

## CloudFront CDN Integration (Optional)

For better performance, use CloudFront:

```go
storage := s3Storage.New(s3Storage.Config{
    Client:     client,
    BucketName: "my-media-bucket",
    Region:     "us-east-1",
    BaseURL:    "https://d111111abcdef8.cloudfront.net", // CloudFront URL
})
```

## Monitoring

Monitor these S3 metrics:
- Upload success/failure rates
- Delete operations
- Storage costs
- Request counts
- Error rates

Use AWS CloudWatch or your monitoring solution of choice.
