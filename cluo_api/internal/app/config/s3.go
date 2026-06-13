package config

// S3Config holds S3-compatible storage configuration (AWS S3 or MinIO).
type S3Config struct {
	Region     string
	BucketName string
	BaseURL    string
	// Endpoint overrides the default AWS endpoint for S3-compatible services (e.g. MinIO).
	// When set, path-style addressing is enabled automatically.
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

func loadS3Config() (S3Config, error) {
	return S3Config{
		Region:          getEnv("AWS_REGION", "us-east-1"),
		BucketName:      getEnv("CLUO_S3_BUCKET_NAME", ""),
		BaseURL:         getEnv("CLUO_S3_BASE_URL", ""),
		Endpoint:        getEnv("CLUO_S3_ENDPOINT", ""),
		AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
	}, nil
}
