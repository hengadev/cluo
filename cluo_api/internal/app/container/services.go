package container

import (
	"context"
	"fmt"

	authService "github.com/hengadev/cluo_api/internal/application/auth"
	caseService "github.com/hengadev/cluo_api/internal/application/case"
	clientService "github.com/hengadev/cluo_api/internal/application/client"
	mediaService "github.com/hengadev/cluo_api/internal/application/media"
	// NOTE: documentService is excluded due to existing compilation errors in the domain layer
	// documentService "github.com/hengadev/cluo_api/internal/application/document"
)

func (c *Container) initServices(ctx context.Context) error {
	// Skip if repositories are not initialized
	if c.caseRepo == nil || c.clientRepo == nil {
		c.logger.WarnContext(ctx, "Repositories not available, skipping service initialization")
		return nil
	}

	// Crypto service is required for services
	if c.crypto == nil {
		c.logger.WarnContext(ctx, "Crypto service not available, skipping service initialization")
		return nil
	}

	// Initialize case service
	c.caseService = caseService.New(c.caseRepo, c.clientRepo, c.caseSubjectRepo, c.crypto)
	c.logger.InfoContext(ctx, "Case service initialized")

	// Initialize client service
	c.clientService = clientService.New(c.clientRepo, c.crypto)
	c.logger.InfoContext(ctx, "Client service initialized")

	// Initialize auth service
	if c.sessionRepo != nil && c.userRepo != nil {
		c.authService = authService.New(c.userRepo, c.sessionRepo, c.crypto)
		c.logger.InfoContext(ctx, "Auth service initialized")
	}

	// NOTE: Document service initialization is commented out due to existing compilation errors
	// in the domain layer. Uncomment when the document domain is fixed.
	// if c.documentRepo != nil && c.documentVersionRepo != nil {
	// 	c.documentService = documentService.New(
	// 		c.documentRepo,
	// 		c.documentVersionRepo,
	// 		c.caseRepo,
	// 		c.clientRepo,
	// 		c.crypto,
	// 	)
	// 	c.logger.InfoContext(ctx, "Document service initialized")
	// }

	// Initialize media service (if storage is available)
	if c.mediaRepo != nil && c.storage != nil {
		c.mediaService = mediaService.New(c.mediaRepo, c.caseRepo, c.storage, c.crypto)
		c.logger.InfoContext(ctx, "Media service initialized")
	}

	return nil
}

// ensureServicesInitialized checks if required services are available.
func (c *Container) ensureServicesInitialized() error {
	if c.caseService == nil {
		return fmt.Errorf("case service not initialized")
	}
	if c.clientService == nil {
		return fmt.Errorf("client service not initialized")
	}
	return nil
}
