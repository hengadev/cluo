package config

// S3Config holds AWS S3 configuration.
type S3Config struct {
	Region     string
	BucketName string
	BaseURL    string
	// AWS credentials (can also use IAM roles)
	AccessKeyID     string
	SecretAccessKey string
}

func loadS3Config() (S3Config, error) {
	return S3Config{
		Region:          getEnv("AWS_REGION", ""),
		BucketName:      getEnv("CLUO_S3_BUCKET_NAME", ""),
		BaseURL:         getEnv("CLUO_S3_BASE_URL", ""),
		AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
	}, nil
}
