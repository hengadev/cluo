package media_test

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

	"github.com/aws/aws-sdk-go-v2/service/s3"
	mediaService "github.com/hengadev/cluo_api/internal/application/media"
	session "github.com/hengadev/cluo_api/internal/common/auth/session"
	services "github.com/hengadev/cluo_api/internal/common/contracts/services"
	envmode "github.com/hengadev/cluo_api/internal/common/envmode"
	logger "github.com/hengadev/cluo_api/internal/common/logger"
	middleware "github.com/hengadev/cluo_api/internal/common/middleware"
	auth "github.com/hengadev/cluo_api/internal/common/middleware/auth"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	caseRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/case"
	mediaRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/media"
	s3Storage "github.com/hengadev/cluo_api/internal/infrastructure/s3"
	mediaHandler "github.com/hengadev/cluo_api/internal/interface/media"
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
	mediaRepo       ports.MediaRepository
	caseRepo        ports.CaseRepository
	vaultSetup      *tu.ServiceVaultSetup
	authCtx         *tu.AuthTestContext
	authSessionRepo session.SessionRepository
	storage         ports.StorageService
	s3Client        *s3.Client
	localstack      *tu.LocalstackContainer
	mediaSvc        ports.MediaService
	handler         mediaHandler.Handler
	testServerURL   string
	testServer      *http.Server
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup logger
	loggerHandler, err := logger.SetHandler("debug", "dev")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	testLogger := slog.New(loggerHandler)
	slog.SetDefault(testLogger)

	// Setup Postgres
	pgContainer, err := tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	// Create pool
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

	// Apply migrations
	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %s\n", err)
	}

	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose migrations: %v", err))
	}
	defer gooseDB.Close()

	if err := goose.UpContext(ctx, gooseDB, "."); err != nil {
		panic(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

	// Setup Redis
	redisContainer, err := tu.SetupRedis(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup redis: %v", err)
	}
	defer tu.TeardownRedis(ctx, nil, redisContainer)

	redisClient = redisContainer.NewClient()

	// Setup Vault
	serviceNames := []string{services.App}
	vaultSetup, err = tu.SetupServiceVault(ctx, nil, serviceNames)
	if err != nil {
		log.Fatalf("Failed to setup service Vault container: %v", err)
	}
	defer tu.TeardownVault(ctx, nil, vaultSetup.VaultContainer)

	// Set environment variables for Vault (for encx library)
	os.Setenv("VAULT_ADDR", vaultSetup.VaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultSetup.VaultContainer.RootToken)

	// Get service-specific crypto service
	var exists bool
	crypto, exists = vaultSetup.GetServiceCrypto(services.App)
	if !exists {
		log.Fatal("App service crypto not found in vault setup")
	}
	if crypto == nil {
		log.Fatal("App crypto service is nil")
	}

	// Setup LocalStack for S3
	localstack, err = tu.SetupLocalstack(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup LocalStack: %v", err)
	}
	defer tu.TeardownLocalstack(ctx, nil, localstack)

	// Create S3 client
	s3Client, err = s3Storage.NewS3Client(ctx, s3Storage.ClientConfig{
		Region:          "us-east-1",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		Endpoint:        localstack.S3Endpoint,
	})
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Create test bucket
	bucketName := "test-media-bucket"
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Fatalf("Failed to create S3 bucket: %v", err)
	}

	// Create storage service
	storage = s3Storage.New(s3Storage.Config{
		Client:     s3Client,
		BucketName: bucketName,
		Region:     "us-east-1",
		BaseURL:    fmt.Sprintf("%s/%s", localstack.S3Endpoint, bucketName),
	})

	// Initialize AuthTestContext for user/session testing
	authCtx = &tu.AuthTestContext{
		Pool:   testPool,
		Redis:  redisClient,
		Crypto: crypto,
	}

	// Create repositories
	mediaRepo = mediaRepository.New(ctx, testPool)
	caseRepo = caseRepository.New(ctx, testPool)

	// Create service
	mediaSvc = mediaService.New(mediaRepo, caseRepo, storage, crypto)

	// Create handler
	authSessionRepo = session.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)
	handler = mediaHandler.New(mediaSvc, authmw)

	// Set required environment variables for logger middleware
	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	// Setup HTTP server
	router := http.NewServeMux()
	handler.RegisterRoutes(router)

	// Use the enhanced AttachLogger middleware
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

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Run tests
	code := m.Run()

	// Graceful shutdown
	log.Println("Shutting down test HTTP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := testServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Test server shutdown failed: %v", err)
	}
	log.Println("Test HTTP server shut down.")

	// Exit with test result code
	os.Exit(code)
}
