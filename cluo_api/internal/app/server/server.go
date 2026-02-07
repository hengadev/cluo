package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/app/container"
	"github.com/hengadev/cluo_api/internal/app/health"
	aiChatHandler "github.com/hengadev/cluo_api/internal/interface/ai_chat"
	aiSpeechToTextHandler "github.com/hengadev/cluo_api/internal/interface/ai_speech_to_text"
	aiTextTransformationHandler "github.com/hengadev/cluo_api/internal/interface/ai_text_transformation"
	aiTranscriptAnalysisHandler "github.com/hengadev/cluo_api/internal/interface/ai_transcript_analysis"
	caseHandler "github.com/hengadev/cluo_api/internal/interface/case"
	clientHandler "github.com/hengadev/cluo_api/internal/interface/client"
	mediaHandler "github.com/hengadev/cluo_api/internal/interface/media"
)

// Server represents the HTTP server.
type Server struct {
	httpServer *http.Server
	container  *container.Container
	config     *config.Config
	logger     *slog.Logger
}

// New creates a new HTTP server.
func New(c *container.Container, cfg *config.Config, logger *slog.Logger) *Server {
	return &Server{
		container: c,
		config:    cfg,
		logger:    logger,
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

	// Register AI routes
	s.registerAIRoutes(mux)
}

func (s *Server) registerCaseRoutes(mux *http.ServeMux) {
	handler := caseHandler.New(s.container.CaseService(), s.container.AuthMiddleware())
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
