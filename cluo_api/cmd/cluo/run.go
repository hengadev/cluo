package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/app/container"
	"github.com/hengadev/cluo_api/internal/app/server"
	"github.com/hengadev/cluo_api/internal/app/tracing"
	"github.com/hengadev/cluo_api/internal/common/envmode"
)

func run(ctx context.Context) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Setup logger
	logger := setupLogger(cfg)

	logger.InfoContext(ctx, "Starting CLUO API",
		"environment", cfg.Environment,
		"address", cfg.Server.Addr(),
	)

	// Setup tracing (if enabled)
	var tracingShutdown func(context.Context) error
	if cfg.Observability.TracingEnabled && cfg.Observability.TracingEndpoint != "" {
		_, shutdown, err := tracing.Setup(ctx, tracing.Config{
			ServiceName:    cfg.Observability.ServiceName,
			ServiceVersion: cfg.Observability.ServiceVersion,
			Environment:    cfg.Environment.String(),
			Endpoint:       cfg.Observability.TracingEndpoint,
			SampleRate:     cfg.Observability.TracingSampler,
		})
		if err != nil {
			logger.WarnContext(ctx, "Failed to setup tracing, continuing without tracing",
				"error", err,
			)
		} else {
			tracingShutdown = shutdown
			logger.InfoContext(ctx, "Tracing initialized",
				"endpoint", cfg.Observability.TracingEndpoint,
			)
		}
	} else {
		tracing.SetupNoop()
		logger.InfoContext(ctx, "Tracing disabled")
	}

	// Create dependency injection container
	ctr, err := container.New(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("create container: %w", err)
	}

	// Create HTTP server
	srv := server.New(ctr, cfg, logger)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start server in goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		if err := srv.Start(ctx); err != nil {
			serverErrCh <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		logger.InfoContext(ctx, "Shutdown signal received")

		// Create shutdown context with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			cfg.Server.ShutdownTimeout,
		)
		defer shutdownCancel()

		// Shutdown server
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.ErrorContext(ctx, "Error shutting down server", "error", err)
		}

		// Shutdown container (database connections, etc.)
		if err := ctr.Shutdown(shutdownCtx); err != nil {
			logger.ErrorContext(ctx, "Error shutting down container", "error", err)
		}

		// Shutdown tracing
		if tracingShutdown != nil {
			if err := tracingShutdown(shutdownCtx); err != nil {
				logger.ErrorContext(ctx, "Error shutting down tracing", "error", err)
			}
		}

		logger.InfoContext(ctx, "Server stopped gracefully")
		return nil

	case err := <-serverErrCh:
		return fmt.Errorf("server error: %w", err)
	}
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	// Determine log level
	var level slog.Level
	switch cfg.Observability.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Determine log style
	switch cfg.Observability.LogStyle {
	case "dev", "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		// Default to JSON in production, text in development
		if cfg.Environment == envmode.Prod {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}
	}

	logger := slog.New(handler)

	// Add service info to all logs
	logger = logger.With(
		"service", cfg.Observability.ServiceName,
		"version", cfg.Observability.ServiceVersion,
	)

	return logger
}
