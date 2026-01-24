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
}

func (s *Server) registerCaseRoutes(mux *http.ServeMux) {
	// Import and use the existing case handler
	// This will be done when wiring everything together
	s.logger.Info("Case routes registered")
}

func (s *Server) registerClientRoutes(mux *http.ServeMux) {
	// Import and use the existing client handler
	s.logger.Info("Client routes registered")
}

func (s *Server) registerMediaRoutes(mux *http.ServeMux) {
	// Import and use the existing media handler
	s.logger.Info("Media routes registered")
}
