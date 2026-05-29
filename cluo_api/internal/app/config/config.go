package config

import (
	"fmt"
	"os"

	"github.com/hengadev/cluo_api/internal/common/envmode"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	Environment   envmode.Mode
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Vault         VaultConfig
	S3            S3Config
	SMTP          SMTPConfig
	Observability ObservabilityConfig
	AI            AIConfig
	RateLimit     RateLimitConfig
}

// Load loads configuration from environment variables.
// It attempts to load a .env file first (for local development).
func Load() (*Config, error) {
	// Try to load .env file - ignore error as it's optional
	_ = godotenv.Load()

	cfg := &Config{}

	// Load environment mode
	envStr := getEnv("CLUO_ENVIRONMENT", "development")
	if err := cfg.Environment.Set(envStr); err != nil {
		return nil, fmt.Errorf("invalid environment: %w", err)
	}

	// Load each config section
	var err error

	cfg.Server, err = loadServerConfig()
	if err != nil {
		return nil, fmt.Errorf("server config: %w", err)
	}

	cfg.Database, err = loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("database config: %w", err)
	}

	cfg.Redis, err = loadRedisConfig()
	if err != nil {
		return nil, fmt.Errorf("redis config: %w", err)
	}

	cfg.Vault, err = loadVaultConfig()
	if err != nil {
		return nil, fmt.Errorf("vault config: %w", err)
	}

	cfg.S3, err = loadS3Config()
	if err != nil {
		return nil, fmt.Errorf("s3 config: %w", err)
	}

	cfg.SMTP = loadSMTPConfig()

	cfg.Observability, err = loadObservabilityConfig(cfg.Environment)
	if err != nil {
		return nil, fmt.Errorf("observability config: %w", err)
	}

	cfg.AI, err = loadAIConfig()
	if err != nil {
		return nil, fmt.Errorf("ai config: %w", err)
	}

	cfg.RateLimit, err = loadRateLimitConfig()
	if err != nil {
		return nil, fmt.Errorf("rate limit config: %w", err)
	}

	// Validate the complete configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return cfg, nil
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvRequired returns the value of an environment variable or an error if not set.
func getEnvRequired(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return value, nil
}
