package portal_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	investigationService "github.com/hengadev/cluo_api/internal/application/investigation"
	rapportService "github.com/hengadev/cluo_api/internal/application/rapport"
	tokenService "github.com/hengadev/cluo_api/internal/application/token"
	session "github.com/hengadev/cluo_api/internal/common/auth/session"
	services "github.com/hengadev/cluo_api/internal/common/contracts/services"
	envmode "github.com/hengadev/cluo_api/internal/common/envmode"
	logger "github.com/hengadev/cluo_api/internal/common/logger"
	middleware "github.com/hengadev/cluo_api/internal/common/middleware"
	auth "github.com/hengadev/cluo_api/internal/common/middleware/auth"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	investigationRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/investigation"
	mediaRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/media"
	rapportRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/rapport"
	tokenRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/token"
	documentRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/document"
	tokenHandler "github.com/hengadev/cluo_api/internal/interface/token"
	migrations "github.com/hengadev/cluo_api/internal/migrations"
	ports "github.com/hengadev/cluo_api/internal/ports"

	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

var (
	testPool        *pgxpool.Pool
	redisClient     *redis.Client
	crypto          encx.CryptoService
	tokenRepo       ports.TokenRepository
	rapportRepo     ports.RapportRepository
	caseRepo        ports.CaseRepository
	vaultSetup      *tu.ServiceVaultSetup
	authCtx         *tu.AuthTestContext
	authSessionRepo session.SessionRepository
	handler         tokenHandler.Handler
	testServerURL   string
	testServer      *http.Server
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggerHandler, err := logger.SetHandler("debug", "dev")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	testLogger := slog.New(loggerHandler)
	slog.SetDefault(testLogger)

	pgContainer, err := tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()

	pgCfg, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config: %v", err))
	}
	pgCfg.MaxConns = 5
	pgCfg.MinConns = 1

	testPool, err = pgxpool.NewWithConfig(poolCtx, pgCfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pool: %v", err))
	}
	if err = testPool.Ping(poolCtx); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %s\n", err)
	}
	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose: %v", err))
	}
	defer gooseDB.Close()
	if err = goose.UpContext(ctx, gooseDB, "."); err != nil {
		panic(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

	redisContainer, err := tu.SetupRedis(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup redis: %v", err)
	}
	defer tu.TeardownRedis(ctx, nil, redisContainer)
	redisClient = redisContainer.NewClient()

	serviceNames := []string{services.App}
	vaultSetup, err = tu.SetupServiceVault(ctx, nil, serviceNames)
	if err != nil {
		log.Fatalf("Failed to setup vault: %v", err)
	}
	defer tu.TeardownVault(ctx, nil, vaultSetup.VaultContainer)

	os.Setenv("VAULT_ADDR", vaultSetup.VaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultSetup.VaultContainer.RootToken)

	var exists bool
	crypto, exists = vaultSetup.GetServiceCrypto(services.App)
	if !exists {
		log.Fatal("App service crypto not found in vault setup")
	}

	authCtx = &tu.AuthTestContext{
		Pool:   testPool,
		Redis:  redisClient,
		Crypto: crypto,
	}

	caseRepo = investigationRepository.New(ctx, testPool)
	tokenRepo = tokenRepository.New(ctx, testPool)
	rapportRepo = rapportRepository.New(ctx, testPool)
	mediaRepo := mediaRepository.New(ctx, testPool)
	documentRepo := documentRepository.New(testPool)

	tokenSvc := tokenService.New(tokenRepo, caseRepo, mediaRepo, crypto)
	rapportSvc := rapportService.New(rapportRepo, caseRepo, crypto)
	caseSvc := investigationService.New(caseRepo, nil, nil, rapportRepo, tokenSvc, crypto)

	authSessionRepo = session.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	// Create a test archive adapter (real repos, fake S3 storage)
	testArchiveAdapter := newTestArchiveAdapter(documentRepo, rapportSvc, mediaRepo, crypto)

	handler = tokenHandler.NewWithArchive(tokenSvc, rapportSvc, documentRepo, crypto, authmw, testArchiveAdapter, caseSvc)

	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	router := http.NewServeMux()
	handler.RegisterRoutes(router)
	loggerMiddleware := middleware.AttachLogger(envmode.Dev, testLogger)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	testServerURL = "http://" + listener.Addr().String()
	testServer = &http.Server{Handler: loggerMiddleware(router)}

	go func() {
		if err := testServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	log.Printf("Test HTTP server started at %s", testServerURL)
	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := testServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Test server shutdown failed: %v", err)
	}

	os.Exit(code)
}
