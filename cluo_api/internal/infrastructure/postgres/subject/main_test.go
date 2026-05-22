package subjectRepository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	subjectRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/subject"
	"github.com/hengadev/cluo_api/internal/migrations"
	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var (
	pgContainer *tu.PostgresContainer
	testPool    *pgxpool.Pool
	repo        ports.CaseSubjectRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error

	// Postgres container
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup postgres container: %v", err))
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	// DB
	log.Println("Creating pgxpool...")
	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()

	config, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse pgxpool config: %v", err))
	}
	config.MaxConns = 5
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	testPool, err = pgxpool.NewWithConfig(poolCtx, config)
	if err != nil {
		tu.TeardownPostgres(ctx, nil, pgContainer)
		panic(fmt.Sprintf("Failed to open test database pool: %v", err))
	}
	log.Println("pgxpool created.")

	if err = testPool.Ping(poolCtx); err != nil {
		panic(fmt.Sprintf("Failed to ping database pool: %v", err))
	}
	log.Println("Database pool ping successful.")

	// Apply migrations
	log.Println("Applying database migrations...")
	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %s\n", err)
	}

	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose migrations: %v", err))
	}
	defer gooseDB.Close()

	if err = goose.UpContext(ctx, gooseDB, "."); err != nil {
		panic(fmt.Sprintf("running all migrations: %s\n", err))
	}
	log.Println("Migrations applied.")

	repo = subjectRepository.New(ctx, testPool)

	// Run tests
	code := m.Run()

	log.Println("Test(s) executed")

	os.Exit(code)
}
