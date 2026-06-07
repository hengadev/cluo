package container

import (
	"context"
	"fmt"

	authService "github.com/hengadev/cluo_api/internal/application/auth"
	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
)

func (c *Container) initAuth(ctx context.Context) error {
	// Skip if Redis client is not available
	if c.redisClient == nil {
		c.logger.WarnContext(ctx, "Redis client not available, skipping auth initialization")
		return nil
	}

	// Crypto service is required for auth
	if c.crypto == nil {
		c.logger.WarnContext(ctx, "Crypto service not available, skipping auth initialization")
		return nil
	}

	// Initialize session repository
	c.sessionRepo = session.NewRedisSessionRepository(c.redisClient)
	c.logger.InfoContext(ctx, "Session repository initialized")

	// Initialize auth service (requires both sessionRepo and userRepo)
	if c.userRepo != nil {
		c.authService = authService.New(c.userRepo, c.sessionRepo, c.crypto)
		c.logger.InfoContext(ctx, "Auth service initialized")
	}

	// Initialize auth middleware
	c.authMiddleware = auth.NewSessionAuthMiddleware(c.sessionRepo, c.crypto, c.vaultClient)
	c.logger.InfoContext(ctx, "Auth middleware initialized")

	return nil
}

// SessionRepository returns the session repository.
func (c *Container) SessionRepository() session.SessionRepository {
	return c.sessionRepo
}

// ensureAuthInitialized checks if auth components are available.
func (c *Container) ensureAuthInitialized() error {
	if c.sessionRepo == nil {
		return fmt.Errorf("session repository not initialized")
	}
	if c.authMiddleware == nil {
		return fmt.Errorf("auth middleware not initialized")
	}
	return nil
}
