package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/app/container"
	"github.com/hengadev/cluo_api/internal/app/health"
	"github.com/hengadev/cluo_api/internal/common/archive"
	mwRatelimit "github.com/hengadev/cluo_api/internal/common/middleware/ratelimit"
	"github.com/hengadev/cluo_api/internal/ports"
	authHandler "github.com/hengadev/cluo_api/internal/interface/auth"
	aiChatHandler "github.com/hengadev/cluo_api/internal/interface/ai_chat"
	aiSpeechToTextHandler "github.com/hengadev/cluo_api/internal/interface/ai_speech_to_text"
	aiTextTransformationHandler "github.com/hengadev/cluo_api/internal/interface/ai_text_transformation"
	aiTranscriptAnalysisHandler "github.com/hengadev/cluo_api/internal/interface/ai_transcript_analysis"
	caseSubjectHandler "github.com/hengadev/cluo_api/internal/interface/case_subject"
	caseTypeHandler "github.com/hengadev/cluo_api/internal/interface/case_type"
	documentHandler "github.com/hengadev/cluo_api/internal/interface/document"
	investigationHandler "github.com/hengadev/cluo_api/internal/interface/investigation"
	clientHandler "github.com/hengadev/cluo_api/internal/interface/client"
	searchHandler "github.com/hengadev/cluo_api/internal/interface/search"
	mediaHandler "github.com/hengadev/cluo_api/internal/interface/media"
	pieceHandler "github.com/hengadev/cluo_api/internal/interface/piece"
	rapportHandler "github.com/hengadev/cluo_api/internal/interface/rapport"
	tokenHandler "github.com/hengadev/cluo_api/internal/interface/token"
)

// Server represents the HTTP server.
type Server struct {
	httpServer      *http.Server
	container       *container.Container
	config          *config.Config
	logger          *slog.Logger
	rateLimitStore  *mwRatelimit.InMemoryStore
}

// New creates a new HTTP server.
func New(c *container.Container, cfg *config.Config, logger *slog.Logger) *Server {
	return &Server{
		container:      c,
		config:         cfg,
		logger:         logger,
		rateLimitStore: mwRatelimit.NewInMemoryStore(),
	}
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	// Create router
	mux := http.NewServeMux()

	// Register routes
	s.registerRoutes(mux)

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.config.Server.Addr(),
		Handler:      mux,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	// Start periodic cleanup of expired rate-limit entries.
	s.rateLimitStore.StartCleanup(ctx, time.Minute)

	s.logger.InfoContext(ctx, "Starting HTTP server",
		"addr", s.config.Server.Addr(),
		"read_timeout", s.config.Server.ReadTimeout,
		"write_timeout", s.config.Server.WriteTimeout,
	)

	// Start server
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.InfoContext(ctx, "Shutting down HTTP server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	s.logger.InfoContext(ctx, "HTTP server stopped")
	return nil
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Register health check routes
	healthChecker := health.NewChecker(
		s.container.DBPool(),
		s.container.RedisClient(),
		s.config.Observability.ServiceVersion,
	)
	healthHandler := health.NewHandler(healthChecker)
	healthHandler.RegisterRoutes(mux)

	// Register API routes (with middleware)
	s.registerAPIRoutes(mux)
}

func (s *Server) registerAPIRoutes(mux *http.ServeMux) {
	// Only register API routes if auth middleware is available
	authMW := s.container.AuthMiddleware()
	if authMW == nil {
		s.logger.Warn("Auth middleware not available, skipping API route registration")
		return
	}

	// Register auth routes
	if s.container.AuthService() != nil {
		s.registerAuthRoutes(mux)
	}

	// Register case routes
	if s.container.CaseService() != nil {
		s.registerCaseRoutes(mux)
	}

	// Register client routes
	if s.container.ClientService() != nil {
		s.registerClientRoutes(mux)
	}

	// Register media routes
	if s.container.MediaService() != nil {
		s.registerMediaRoutes(mux)
	}

	// Register piece routes
	if s.container.PieceService() != nil {
		s.registerPieceRoutes(mux)
	}

	// Register rapport routes
	if s.container.RapportService() != nil {
		s.registerRapportRoutes(mux)
	}

	// Register token routes
	if s.container.TokenService() != nil {
		s.registerTokenRoutes(mux)
	}

	// Register case type routes
	if s.container.CaseTypeService() != nil {
		s.registerCaseTypeRoutes(mux)
	}

	// Register subject routes
	if s.container.CaseSubjectService() != nil {
		s.registerSubjectRoutes(mux)
	}

	// Register document routes
	if s.container.DocumentService() != nil {
		s.registerDocumentRoutes(mux)
	}

	// Register search routes
	if s.container.CaseService() != nil && s.container.ClientService() != nil {
		s.registerSearchRoutes(mux)
	}

	// Register AI routes
	s.registerAIRoutes(mux)
}

func (s *Server) registerCaseRoutes(mux *http.ServeMux) {
	handler := investigationHandler.New(s.container.CaseService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Case routes registered")
}

func (s *Server) registerClientRoutes(mux *http.ServeMux) {
	handler := clientHandler.New(s.container.ClientService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Client routes registered")
}

func (s *Server) registerMediaRoutes(mux *http.ServeMux) {
	handler := mediaHandler.New(s.container.MediaService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Media routes registered")
}

func (s *Server) registerPieceRoutes(mux *http.ServeMux) {
	handler := pieceHandler.New(s.container.PieceService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Piece routes registered")
}

func (s *Server) registerRapportRoutes(mux *http.ServeMux) {
	handler := rapportHandler.New(s.container.RapportService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Rapport routes registered")
}

func (s *Server) registerTokenRoutes(mux *http.ServeMux) {
	var handler tokenHandler.Handler

	// If storage is available, create the archive adapter for download support.
	if s.container.StorageService() != nil {
		archiveDeps := archive.NewAdapter(
			s.container.TypedDocumentRepository(),
			s.container.RapportService(),
			s.container.MediaRepository().(ports.MediaRepository),
			s.container.StorageService(),
			s.container.Crypto(),
		)
		handler = tokenHandler.NewWithArchive(
			s.container.TokenService(),
			s.container.RapportService(),
			s.container.TypedDocumentRepository(),
			s.container.Crypto(),
			s.container.AuthMiddleware(),
			archiveDeps,
			s.container.CaseService(),
		)
	} else {
		handler = tokenHandler.New(
			s.container.TokenService(),
			s.container.RapportService(),
			s.container.TypedDocumentRepository(),
			s.container.Crypto(),
			s.container.AuthMiddleware(),
			s.container.CaseService(),
		)
	}

	// Apply rate limiter to the public portal token resolution route.
	maxReqs, window := s.config.RateLimit.TokenWindow()
	handler = handler.(*tokenHandler.TokenHandler).WithTokenRateLimiter(
		mwRatelimit.RateLimiter(s.rateLimitStore, mwRatelimit.Config{MaxRequests: maxReqs, Window: window}),
	)

	handler.RegisterRoutes(mux)
	s.logger.Info("Token routes registered")
}

func (s *Server) registerCaseTypeRoutes(mux *http.ServeMux) {
	handler := caseTypeHandler.New(s.container.CaseTypeService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("CaseType routes registered")
}

func (s *Server) registerSubjectRoutes(mux *http.ServeMux) {
	handler := caseSubjectHandler.New(s.container.CaseSubjectService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Subject routes registered")
}

func (s *Server) registerAIRoutes(mux *http.ServeMux) {
	authMW := s.container.AuthMiddleware()
	if authMW == nil {
		s.logger.Warn("Auth middleware not available, skipping AI route registration")
		return
	}

	// Register text transformation routes
	if s.container.TextTransformationService() != nil {
		handler := aiTextTransformationHandler.New(s.container.TextTransformationService(), authMW)
		handler.RegisterRoutes(mux)
		s.logger.Info("Text transformation routes registered")
	}

	// Register speech-to-text routes
	if s.container.SpeechToTextService() != nil {
		handler := aiSpeechToTextHandler.New(s.container.SpeechToTextService(), authMW)
		handler.RegisterRoutes(mux)
		s.logger.Info("Speech-to-text routes registered")
	}

	// Register transcript analysis routes
	if s.container.TranscriptAnalysisService() != nil {
		handler := aiTranscriptAnalysisHandler.New(s.container.TranscriptAnalysisService(), authMW)
		handler.RegisterRoutes(mux)
		s.logger.Info("Transcript analysis routes registered")
	}

	// Register chat routes
	if s.container.ChatService() != nil {
		handler := aiChatHandler.New(s.container.ChatService(), authMW)
		handler.RegisterRoutes(mux)
		s.logger.Info("Chat routes registered")
	}
}

func (s *Server) registerAuthRoutes(mux *http.ServeMux) {
	handler := authHandler.New(s.container.AuthService(), s.container.AuthMiddleware())

	// Apply rate limiter to login route.
	maxReqs, window := s.config.RateLimit.LoginWindow()
	handler = handler.(*authHandler.AuthHandler).WithLoginRateLimiter(
		mwRatelimit.RateLimiter(s.rateLimitStore, mwRatelimit.Config{MaxRequests: maxReqs, Window: window}),
	)

	handler.RegisterRoutes(mux)
	s.logger.Info("Auth routes registered")
}

func (s *Server) registerDocumentRoutes(mux *http.ServeMux) {
	handler := documentHandler.New(s.container.DocumentService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Document routes registered")
}

func (s *Server) registerSearchRoutes(mux *http.ServeMux) {
	handler := searchHandler.New(s.container.SearchService(), s.container.AuthMiddleware())
	handler.RegisterRoutes(mux)
	s.logger.Info("Search routes registered")
}
