package container

import (
	"context"
	"fmt"

	caseRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/case"
	caseSubjectRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/case_subject"
	clientRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/client"
	mediaRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/media"
	// NOTE: documentRepository is excluded due to existing compilation errors
	// documentRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/document"
)

func (c *Container) initRepositories(ctx context.Context) error {
	// Skip if no database pool
	if c.dbPool == nil {
		c.logger.WarnContext(ctx, "Database pool not available, skipping repository initialization")
		return nil
	}

	// Initialize case repository
	c.caseRepo = caseRepository.New(ctx, c.dbPool)
	c.logger.InfoContext(ctx, "Case repository initialized")

	// Initialize case subject repository
	c.caseSubjectRepo = caseSubjectRepository.New(ctx, c.dbPool)
	c.logger.InfoContext(ctx, "Case subject repository initialized")

	// Initialize client repository
	c.clientRepo = clientRepository.New(ctx, c.dbPool)
	c.logger.InfoContext(ctx, "Client repository initialized")

	// NOTE: Document repositories are commented out due to existing compilation errors
	// c.documentRepo = documentRepository.New(c.dbPool)
	// c.logger.InfoContext(ctx, "Document repository initialized")
	// c.documentVersionRepo = documentRepository.NewVersionRepository(c.dbPool)
	// c.logger.InfoContext(ctx, "Document version repository initialized")

	// Initialize media repository
	c.mediaRepo = mediaRepository.New(ctx, c.dbPool)
	c.logger.InfoContext(ctx, "Media repository initialized")

	return nil
}

// CaseRepository returns the case repository.
func (c *Container) CaseRepository() interface{} {
	return c.caseRepo
}

// ClientRepository returns the client repository.
func (c *Container) ClientRepository() interface{} {
	return c.clientRepo
}

// DocumentRepository returns the document repository.
func (c *Container) DocumentRepository() interface{} {
	return c.documentRepo
}

// MediaRepository returns the media repository.
func (c *Container) MediaRepository() interface{} {
	return c.mediaRepo
}

// ensureRepositoriesInitialized checks if required repositories are available.
func (c *Container) ensureRepositoriesInitialized() error {
	if c.caseRepo == nil {
		return fmt.Errorf("case repository not initialized")
	}
	if c.clientRepo == nil {
		return fmt.Errorf("client repository not initialized")
	}
	return nil
}
