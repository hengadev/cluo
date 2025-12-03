package clientHelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"
	clientRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/client"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearClientsTable truncates the contacts table for clean test state
func ClearClientsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.contacts RESTART IDENTITY CASCADE", clientRepository.Schema))
	require.NoError(t, err)
}

// NewTestContact creates a Contact domain object with basic test data (plaintext fields only)
func NewTestClient(t *testing.T) *client.Client {
	t.Helper()

	return &client.Client{
		ID:         uuid.New(),
		Name:       "client_name",
		Type:       client.ClientTypePerson,
		ContactIDs: []string{},
		CreatedAt:  time.Now(),
	}
}

// NewTestClientEncx creates a mock ClientEncx domain object with basic test data (plaintext fields only)
func NewTestClientEncx(t *testing.T) *client.ClientEncx {
	t.Helper()
	return &client.ClientEncx{
		ID:                  uuid.New(),
		NameEncrypted:       []byte("name_encrypted"),
		NameHash:            "name_hash",
		TypeEncrypted:       []byte("type_encrypted"),
		TypeHash:            "type_hash",
		ContactIDsEncrypted: []byte("contact_ids_encrypted"),
		CreatedAt:           time.Now(),
		DEKEncrypted:        []byte("dek_encrypted"),
		KeyVersion:          1,
	}
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
