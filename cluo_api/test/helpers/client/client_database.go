package clientHelpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
	clientRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/client"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearClientsTable truncates the contacts table for clean test state
func ClearClientsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.clients RESTART IDENTITY CASCADE", clientRepository.Schema))
	require.NoError(t, err)
}

// InsertClientEncx inserts a ClientEncx record into the database for testing
func InsertClientEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientEncx *client.ClientEncx) error {
	t.Helper()

	query := fmt.Sprintf(`
		INSERT INTO %s.clients (
			id, created_at, name_encrypted, name_hash, type_encrypted, type_hash,
			contactids_encrypted, dek_encrypted, key_version, metadata
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, clientRepository.Schema)

	_, err := pool.Exec(ctx, query,
		clientEncx.ID, clientEncx.CreatedAt, clientEncx.NameEncrypted, clientEncx.NameHash,
		clientEncx.TypeEncrypted, clientEncx.TypeHash, clientEncx.ContactIDsEncrypted,
		clientEncx.DEKEncrypted, clientEncx.KeyVersion, clientEncx.Metadata,
	)

	return err
}

// GetClientEncxByID retrieves a client by ID from the database for testing
func GetClientEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientID uuid.UUID) (*client.ClientEncx, error) {
	t.Helper()

	query := fmt.Sprintf(`
		SELECT
			id, created_at, name_encrypted, name_hash, type_encrypted, type_hash,
			contactids_encrypted, dek_encrypted, key_version, metadata
		FROM %s.clients WHERE id = $1
	`, clientRepository.Schema)

	clientEncx := &client.ClientEncx{}

	err := pool.QueryRow(ctx, query, clientID).Scan(
		&clientEncx.ID, &clientEncx.CreatedAt, &clientEncx.NameEncrypted, &clientEncx.NameHash,
		&clientEncx.TypeEncrypted, &clientEncx.TypeHash, &clientEncx.ContactIDsEncrypted,
		&clientEncx.DEKEncrypted, &clientEncx.KeyVersion, &clientEncx.Metadata,
	)

	return clientEncx, err
}
