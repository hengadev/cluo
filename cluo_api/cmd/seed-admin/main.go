package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	"github.com/hengadev/cluo_api/internal/domain/user"
	userRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/user"
	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("seed-admin: %v\n", err)
	}
}

func run(ctx context.Context) error {
	_ = godotenv.Load()

	email := flag.String("email", getEnv("SEED_ADMIN_EMAIL", ""), "admin email address")
	password := flag.String("password", getEnv("SEED_ADMIN_PASSWORD", ""), "admin password (min 8 chars)")
	flag.Parse()

	if *email == "" {
		return fmt.Errorf("email is required: use --email or SEED_ADMIN_EMAIL")
	}
	if *password == "" {
		return fmt.Errorf("password is required: use --password or SEED_ADMIN_PASSWORD")
	}
	if len(*password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	pool, err := connectDB(ctx)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer pool.Close()

	crypto, err := initCrypto(ctx)
	if err != nil {
		return fmt.Errorf("crypto: %w", err)
	}

	userRepo := userRepository.New(ctx, pool)

	emailBytes, err := encx.SerializeValue(*email)
	if err != nil {
		return fmt.Errorf("serialize email: %w", err)
	}
	emailHash := crypto.HashBasic(ctx, emailBytes)

	exists, err := userRepo.ExistsByEmailHash(ctx, emailHash)
	if err != nil {
		return fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		log.Printf("admin user already exists, skipping")
		return nil
	}

	newUser := &user.User{
		ID:        uuid.New(),
		Email:     *email,
		Password:  *password,
		Role:      identity.Administrator.String(),
		CreatedAt: time.Now(),
	}

	userEncx, err := user.ProcessUserEncx(ctx, crypto, newUser)
	if err != nil {
		return fmt.Errorf("encrypt user: %w", err)
	}

	if err := userRepo.CreateUser(ctx, userEncx); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	log.Printf("admin user created successfully")
	return nil
}

func connectDB(ctx context.Context) (*pgxpool.Pool, error) {
	host := os.Getenv("CLUO_DB_HOST")
	name := os.Getenv("CLUO_DB_NAME")
	u := os.Getenv("CLUO_DB_USER")
	pass := os.Getenv("CLUO_DB_PASSWORD")
	if host == "" || name == "" || u == "" || pass == "" {
		return nil, fmt.Errorf("CLUO_DB_HOST, CLUO_DB_NAME, CLUO_DB_USER, CLUO_DB_PASSWORD are all required")
	}
	port := getEnv("CLUO_DB_PORT", "5432")
	sslMode := getEnv("CLUO_DB_SSL_MODE", "require")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		u, url.QueryEscape(pass), host, port, name, sslMode)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

func initCrypto(ctx context.Context) (encx.CryptoService, error) {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		return nil, fmt.Errorf("VAULT_ADDR is required")
	}

	vaultCfg := api.DefaultConfig()
	vaultCfg.Address = vaultAddr

	client, err := api.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("create vault client: %w", err)
	}

	roleID := os.Getenv("VAULT_APPROLE_ROLE_ID")
	secretID := os.Getenv("VAULT_APPROLE_SECRET_ID")
	if roleID != "" && secretID != "" {
		resp, err := client.Logical().Write("auth/approle/login", map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		})
		if err != nil {
			return nil, fmt.Errorf("approle login: %w", err)
		}
		client.SetToken(resp.Auth.ClientToken)
		// encx providers read VAULT_TOKEN from the environment
		os.Setenv("VAULT_TOKEN", resp.Auth.ClientToken)
	} else if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	} else {
		return nil, fmt.Errorf("vault auth required: set VAULT_APPROLE_ROLE_ID/VAULT_APPROLE_SECRET_ID or VAULT_TOKEN")
	}

	kms, err := hashicorpkeys.NewTransitService()
	if err != nil {
		return nil, fmt.Errorf("create KMS provider: %w", err)
	}

	secrets, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		return nil, fmt.Errorf("create secrets store: %w", err)
	}

	cfg := encx.Config{
		KEKAlias:    "cluo-encryption-key",
		PepperAlias: "cluo",
		DBPath:      "data/.encx",
	}

	crypto, err := encx.NewCrypto(ctx, kms, secrets, cfg)
	if err != nil {
		return nil, fmt.Errorf("create crypto service: %w", err)
	}

	return crypto, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
