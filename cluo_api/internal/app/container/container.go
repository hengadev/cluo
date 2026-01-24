package container

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/vault/api"
	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Container holds all application dependencies.
type Container struct {
	config *config.Config
	logger *slog.Logger

	// Infrastructure clients
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	vaultClient *api.Client
	crypto      encx.CryptoService

	// Repositories
	caseRepo            ports.CaseRepository
	caseSubjectRepo     ports.CaseSubjectRepository
	clientRepo          ports.ClientRepository
	documentRepo        ports.DocumentRepository
	documentVersionRepo ports.DocumentVersionRepository
	mediaRepo           ports.MediaRepository

	// Services
	caseService     ports.CaseService
	clientService   ports.ClientService
	documentService ports.DocumentService
	mediaService    ports.MediaService
	storage         ports.StorageService

	// Auth
	sessionRepo    session.SessionRepository
	authMiddleware auth.AuthMiddleware
}

// New creates a new dependency injection container.
func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Container, error) {
	c := &Container{
		config: cfg,
		logger: logger,
	}

	// Initialize infrastructure clients
	if err := c.initInfrastructure(ctx); err != nil {
		c.Shutdown(ctx) // Clean up any initialized resources
		return nil, fmt.Errorf("init infrastructure: %w", err)
	}

	// Initialize repositories
	if err := c.initRepositories(ctx); err != nil {
		c.Shutdown(ctx)
		return nil, fmt.Errorf("init repositories: %w", err)
	}

	// Initialize services
	if err := c.initServices(ctx); err != nil {
		c.Shutdown(ctx)
		return nil, fmt.Errorf("init services: %w", err)
	}

	// Initialize auth middleware
	if err := c.initAuth(ctx); err != nil {
		c.Shutdown(ctx)
		return nil, fmt.Errorf("init auth: %w", err)
	}

	logger.InfoContext(ctx, "Container initialized successfully")
	return c, nil
}

// Shutdown gracefully shuts down all container resources.
func (c *Container) Shutdown(ctx context.Context) error {
	var errs []error

	if c.dbPool != nil {
		c.dbPool.Close()
		c.logger.InfoContext(ctx, "Database connection pool closed")
	}

	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		} else {
			c.logger.InfoContext(ctx, "Redis connection closed")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// Config returns the application configuration.
func (c *Container) Config() *config.Config {
	return c.config
}

// Logger returns the application logger.
func (c *Container) Logger() *slog.Logger {
	return c.logger
}

// DBPool returns the database connection pool.
func (c *Container) DBPool() *pgxpool.Pool {
	return c.dbPool
}

// RedisClient returns the Redis client.
func (c *Container) RedisClient() *redis.Client {
	return c.redisClient
}

// VaultClient returns the Vault client.
func (c *Container) VaultClient() *api.Client {
	return c.vaultClient
}

// Crypto returns the crypto service.
func (c *Container) Crypto() encx.CryptoService {
	return c.crypto
}

// CaseService returns the case service.
func (c *Container) CaseService() ports.CaseService {
	return c.caseService
}

// ClientService returns the client service.
func (c *Container) ClientService() ports.ClientService {
	return c.clientService
}

// DocumentService returns the document service.
func (c *Container) DocumentService() ports.DocumentService {
	return c.documentService
}

// MediaService returns the media service.
func (c *Container) MediaService() ports.MediaService {
	return c.mediaService
}

// AuthMiddleware returns the authentication middleware.
func (c *Container) AuthMiddleware() auth.AuthMiddleware {
	return c.authMiddleware
}
