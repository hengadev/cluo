package container

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/vault/api"
	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/infrastructure/worker"
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
	caseTypeRepo        ports.CaseTypeRepository
	clientRepo          ports.ClientRepository
	documentRepo        ports.DocumentRepository
	documentVersionRepo ports.DocumentVersionRepository
	mediaRepo           ports.MediaRepository
	pieceRepo           ports.PieceRepository
	userRepo            ports.UserRepository
	rapportRepo         ports.RapportRepository
	tokenRepo           ports.TokenRepository

	// Services
	caseService        ports.CaseService
	caseSubjectService ports.CaseSubjectService
	caseTypeService    ports.CaseTypeService
	clientService      ports.ClientService
	searchService      ports.SearchService
	documentService    ports.DocumentService
	mediaService       ports.MediaService
	pieceService       ports.PieceService
	storage            ports.StorageService
	authService        ports.AuthService
	rapportService     ports.RapportService
	tokenService       ports.TokenService

	// AI Services
	textTransformationService ports.TextTransformationService
	speechToTextService       ports.SpeechToTextService
	transcriptAnalysisService ports.TranscriptAnalysisService
	transcriptionRepo         ports.TranscriptionRepository
	transcriptionJobRepo      ports.TranscriptionJobRepository
	transcriptionWorker       *worker.TranscriptionWorker
	chatService              ports.ChatService

	// Auth
	sessionRepo    *session.RedisSessionRepository
	authMiddleware auth.AuthMiddleware

	// Email
	emailService ports.EmailService
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

	// Initialize AI services
	if err := c.initAIServices(ctx); err != nil {
		// AI services are optional, log warning but continue
		c.logger.WarnContext(ctx, "AI services initialization failed", "error", err)
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

// SearchService returns the search service.
func (c *Container) SearchService() ports.SearchService {
	return c.searchService
}

// DocumentService returns the document service.
func (c *Container) DocumentService() ports.DocumentService {
	return c.documentService
}

// MediaService returns the media service.
func (c *Container) MediaService() ports.MediaService {
	return c.mediaService
}

// PieceService returns the piece service.
func (c *Container) PieceService() ports.PieceService {
	return c.pieceService
}

// AuthService returns the auth service.
func (c *Container) AuthService() ports.AuthService {
	return c.authService
}

// AuthMiddleware returns the authentication middleware.
func (c *Container) AuthMiddleware() auth.AuthMiddleware {
	return c.authMiddleware
}

// TextTransformationService returns the text transformation service.
func (c *Container) TextTransformationService() ports.TextTransformationService {
	return c.textTransformationService
}

// SpeechToTextService returns the speech-to-text service.
func (c *Container) SpeechToTextService() ports.SpeechToTextService {
	return c.speechToTextService
}

// TranscriptAnalysisService returns the transcript analysis service.
func (c *Container) TranscriptAnalysisService() ports.TranscriptAnalysisService {
	return c.transcriptAnalysisService
}

// ChatService returns the chat service.
func (c *Container) ChatService() ports.ChatService {
	return c.chatService
}

// RapportService returns the rapport service.
func (c *Container) RapportService() ports.RapportService {
	return c.rapportService
}

// RapportRepository returns the rapport repository.
func (c *Container) RapportRepository() ports.RapportRepository {
	return c.rapportRepo
}

// TokenService returns the token service.
func (c *Container) TokenService() ports.TokenService {
	return c.tokenService
}

// CaseTypeService returns the case type service.
func (c *Container) CaseTypeService() ports.CaseTypeService {
	return c.caseTypeService
}

// CaseSubjectService returns the case subject service.
func (c *Container) CaseSubjectService() ports.CaseSubjectService {
	return c.caseSubjectService
}

// TypedDocumentRepository returns the document repository with its correct type.
func (c *Container) TypedDocumentRepository() ports.DocumentRepository {
	return c.documentRepo
}

// StorageService returns the storage service.
func (c *Container) StorageService() ports.StorageService {
	return c.storage
}

// EmailService returns the email service.
func (c *Container) EmailService() ports.EmailService {
	return c.emailService
}

// StartBackgroundWorkers starts all background workers.
func (c *Container) StartBackgroundWorkers(ctx context.Context) {
	if c.transcriptionWorker != nil {
		c.transcriptionWorker.Start(ctx)
		c.logger.InfoContext(ctx, "Transcription worker started")
	}
}

// StopBackgroundWorkers stops all background workers gracefully.
func (c *Container) StopBackgroundWorkers() {
	if c.transcriptionWorker != nil {
		c.transcriptionWorker.Stop()
		c.logger.Info("Transcription worker stopped")
	}
}

