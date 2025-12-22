package main

import "context"

func run(ctx context.Context) error {
	// // Setup logger
	// logger := setupLogger()
	//
	// // Load configuration
	// cfg, err := app.LoadConfig(ctx)
	// if err != nil {
	// 	return fmt.Errorf("load config: %w", err)
	// }
	//
	// logger.InfoContext(ctx, "Starting Leviosa modular monolith",
	// 	"environment", cfg.Environment,
	// 	"port", cfg.ServerPort,
	// )
	//
	// // Create dependency injection container
	// container, err := app.NewContainer(ctx, cfg)
	// if err != nil {
	// 	return fmt.Errorf("create container: %w", err)
	// }
	//
	// // Create HTTP server
	// server := app.NewServer(container, logger)
	//
	// // Setup graceful shutdown
	// ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// defer cancel()
	//
	// // Start server in goroutine
	// serverErrCh := make(chan error, 1)
	// go func() {
	// 	if err := server.Start(ctx); err != nil {
	// 		serverErrCh <- err
	// 	}
	// }()
	//
	// // Wait for shutdown signal or server error
	// select {
	// case <-ctx.Done():
	// 	logger.InfoContext(ctx, "Shutdown signal received")
	//
	// 	// Give server 30 seconds to shutdown gracefully
	// 	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	// 	defer shutdownCancel()
	//
	// 	if err := server.Shutdown(shutdownCtx); err != nil {
	// 		return fmt.Errorf("shutdown server: %w", err)
	// 	}
	//
	// 	logger.InfoContext(ctx, "Server stopped gracefully")
	// 	return nil
	//
	// case err := <-serverErrCh:
	// 	return fmt.Errorf("server error: %w", err)
	// }

	return nil
}
