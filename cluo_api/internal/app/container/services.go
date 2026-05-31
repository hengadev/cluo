package container

import (
	"context"
	"fmt"

	authService "github.com/hengadev/cluo_api/internal/application/auth"
	caseSubjectService "github.com/hengadev/cluo_api/internal/application/case_subject"
	caseTypeService "github.com/hengadev/cluo_api/internal/application/case_type"
	investigationService "github.com/hengadev/cluo_api/internal/application/investigation"
	clientService "github.com/hengadev/cluo_api/internal/application/client"
	mediaService "github.com/hengadev/cluo_api/internal/application/media"
	pieceService "github.com/hengadev/cluo_api/internal/application/piece"
	rapportService "github.com/hengadev/cluo_api/internal/application/rapport"
	searchService "github.com/hengadev/cluo_api/internal/application/search"
	tokenService "github.com/hengadev/cluo_api/internal/application/token"
	documentService "github.com/hengadev/cluo_api/internal/application/document"
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

	// Initialize token service early so caseService can depend on it.
	if c.tokenRepo != nil && c.mediaRepo != nil {
		c.tokenService = tokenService.New(c.tokenRepo, c.caseRepo, c.mediaRepo, c.clientRepo, c.crypto, c.emailService, c.config.SMTP.PortalPublicURL, c.logger)
		c.logger.InfoContext(ctx, "Token service initialized")
	}

	// Initialize case service
	c.caseService = investigationService.New(c.caseRepo, c.clientRepo, c.caseSubjectRepo, c.rapportRepo, c.tokenService, c.crypto)
	c.logger.InfoContext(ctx, "Case service initialized")

	// Initialize client service
	c.clientService = clientService.New(c.clientRepo, c.crypto)
	c.logger.InfoContext(ctx, "Client service initialized")

	// Initialize search service
	c.searchService = searchService.New(c.caseService, c.clientService)
	c.logger.InfoContext(ctx, "Search service initialized")

	// Initialize auth service
	if c.sessionRepo != nil && c.userRepo != nil {
		c.authService = authService.New(c.userRepo, c.sessionRepo, c.crypto)
		c.logger.InfoContext(ctx, "Auth service initialized")
	}

	if c.documentRepo != nil && c.documentVersionRepo != nil {
		c.documentService = documentService.New(
			c.documentRepo,
			c.documentVersionRepo,
			c.caseRepo,
			c.clientRepo,
			c.crypto,
			c.emailService,
			c.logger,
		)
		c.logger.InfoContext(ctx, "Document service initialized")
	}

	// Initialize media service (if storage is available)
	if c.mediaRepo != nil && c.storage != nil {
		c.mediaService = mediaService.New(c.mediaRepo, c.caseRepo, c.storage, c.crypto)
		c.logger.InfoContext(ctx, "Media service initialized")
	}

	// Initialize piece service (if storage is available)
	if c.pieceRepo != nil && c.storage != nil {
		c.pieceService = pieceService.New(c.pieceRepo, c.caseRepo, c.storage, c.crypto)
		c.logger.InfoContext(ctx, "Piece service initialized")
	}

	// Initialize rapport service
	if c.rapportRepo != nil {
		c.rapportService = rapportService.New(c.rapportRepo, c.caseRepo, c.crypto)
		c.logger.InfoContext(ctx, "Rapport service initialized")
	}

	// Initialize case type service
	if c.caseTypeRepo != nil {
		c.caseTypeService = caseTypeService.New(c.caseTypeRepo)
		c.logger.InfoContext(ctx, "CaseType service initialized")
	}

	// Initialize case subject service
	if c.caseSubjectRepo != nil {
		c.caseSubjectService = caseSubjectService.New(c.caseSubjectRepo, c.crypto)
		c.logger.InfoContext(ctx, "CaseSubject service initialized")
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
