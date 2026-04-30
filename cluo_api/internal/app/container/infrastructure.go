package container

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/cluo_api/internal/common/envmode"
	s3Storage "github.com/hengadev/cluo_api/internal/infrastructure/s3"
	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func (c *Container) initInfrastructure(ctx context.Context) error {
	// Initialize database pool
	if err := c.initDatabase(ctx); err != nil {
		return fmt.Errorf("database: %w", err)
	}

	// Initialize Redis client
	if err := c.initRedis(ctx); err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	// Initialize Vault client (optional in development)
	if err := c.initVault(ctx); err != nil {
		// In development, Vault is optional
		if c.config.Environment != envmode.Dev {
			return fmt.Errorf("vault: %w", err)
		}
		c.logger.WarnContext(ctx, "Vault initialization skipped in development mode", "error", err)
	}

	// Initialize crypto service
	if err := c.initCrypto(ctx); err != nil {
		if c.config.Environment != envmode.Dev {
			return fmt.Errorf("crypto: %w", err)
		}
		c.logger.WarnContext(ctx, "Crypto service initialization skipped in development mode", "error", err)
	}

	// Initialize S3 storage
	if err := c.initStorage(ctx); err != nil {
		if c.config.Environment != envmode.Dev {
			return fmt.Errorf("storage: %w", err)
		}
		c.logger.WarnContext(ctx, "S3 storage initialization skipped in development mode", "error", err)
	}

	return nil
}

func (c *Container) initDatabase(ctx context.Context) error {
	cfg := c.config.Database

	// Skip if no database host configured (development without DB)
	if cfg.Host == "" {
		c.logger.WarnContext(ctx, "Database host not configured, skipping database initialization")
		return nil
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return fmt.Errorf("parse connection string: %w", err)
	}

	// Apply pool configuration
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	c.dbPool = pool
	c.logger.InfoContext(ctx, "Database connection pool initialized",
		"host", cfg.Host,
		"database", cfg.Name,
		"max_conns", cfg.MaxOpenConns,
	)

	return nil
}

func (c *Container) initRedis(ctx context.Context) error {
	cfg := c.config.Redis

	// Skip if no Redis host configured
	if cfg.Host == "" {
		c.logger.WarnContext(ctx, "Redis host not configured, skipping Redis initialization")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return fmt.Errorf("ping redis: %w", err)
	}

	c.redisClient = client
	c.logger.InfoContext(ctx, "Redis connection initialized", "addr", cfg.Addr())

	return nil
}

func (c *Container) initVault(ctx context.Context) error {
	cfg := c.config.Vault

	// Skip if no Vault address configured
	if cfg.Address == "" {
		return fmt.Errorf("vault address not configured")
	}

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.Address

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}

	// Set authentication
	if cfg.UseAppRole() {
		// AppRole authentication
		data := map[string]interface{}{
			"role_id":   cfg.AppRoleID,
			"secret_id": cfg.AppRoleSecretID,
		}
		resp, err := client.Logical().Write("auth/approle/login", data)
		if err != nil {
			return fmt.Errorf("approle login: %w", err)
		}
		client.SetToken(resp.Auth.ClientToken)
	} else if cfg.Token != "" {
		// Token authentication (development)
		client.SetToken(cfg.Token)
	} else {
		return fmt.Errorf("no vault authentication method configured")
	}

	// Test connection
	_, err = client.Sys().Health()
	if err != nil {
		return fmt.Errorf("vault health check: %w", err)
	}

	c.vaultClient = client
	c.logger.InfoContext(ctx, "Vault client initialized", "address", cfg.Address)

	return nil
}

func (c *Container) initCrypto(ctx context.Context) error {
	// Vault client must be initialized first
	if c.vaultClient == nil {
		return fmt.Errorf("vault client not initialized")
	}

	// TODO: os.Setenv is not thread-safe in Go. The encx hashicorp providers currently
	// require VAULT_ADDR and VAULT_TOKEN to be set as environment variables during
	// initialization. This should be refactored to pass config directly once encx
	// supports explicit configuration instead of reading from env vars.
	originalAddr := os.Getenv("VAULT_ADDR")
	originalToken := os.Getenv("VAULT_TOKEN")

	os.Setenv("VAULT_ADDR", c.config.Vault.Address)
	os.Setenv("VAULT_TOKEN", c.vaultClient.Token())

	// Create KMS provider
	kms, err := hashicorpkeys.NewTransitService()
	if err != nil {
		os.Setenv("VAULT_ADDR", originalAddr)
		os.Setenv("VAULT_TOKEN", originalToken)
		return fmt.Errorf("create KMS provider: %w", err)
	}

	// Create secrets provider
	secrets, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		os.Setenv("VAULT_ADDR", originalAddr)
		os.Setenv("VAULT_TOKEN", originalToken)
		return fmt.Errorf("create secrets store: %w", err)
	}

	// Restore original environment variables
	os.Setenv("VAULT_ADDR", originalAddr)
	os.Setenv("VAULT_TOKEN", originalToken)

	// Create crypto config
	cfg := encx.Config{
		KEKAlias:    "cluo-encryption-key",
		PepperAlias: "cluo",
		DBPath:      "data/.encx",
	}

	// Create crypto service
	crypto, err := encx.NewCrypto(ctx, kms, secrets, cfg)
	if err != nil {
		return fmt.Errorf("create crypto service: %w", err)
	}

	c.crypto = crypto
	c.logger.InfoContext(ctx, "Crypto service initialized")

	return nil
}

func (c *Container) initStorage(ctx context.Context) error {
	cfg := c.config.S3

	// Skip if no S3 bucket configured
	if cfg.BucketName == "" {
		c.logger.WarnContext(ctx, "S3 bucket not configured, skipping S3 initialization")
		return nil
	}

	// Create AWS config
	var awsCfg aws.Config
	var err error

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		// Use explicit credentials
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.Region),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			)),
		)
	} else {
		// Use default credential chain (IAM role, env vars, etc.)
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.Region),
		)
	}

	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsCfg)

	// Create storage service
	c.storage = s3Storage.New(s3Storage.Config{
		Client:     s3Client,
		BucketName: cfg.BucketName,
		Region:     cfg.Region,
		BaseURL:    cfg.BaseURL,
	})

	c.logger.InfoContext(ctx, "S3 storage initialized",
		"bucket", cfg.BucketName,
		"region", cfg.Region,
	)

	return nil
}
