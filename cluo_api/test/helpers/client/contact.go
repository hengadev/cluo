package clientHelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/hengadev/cluo_api/internal/infrastructure/postgres/client"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearContactsTable truncates the contacts table for clean test state
func ClearContactsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.contacts RESTART IDENTITY CASCADE", clientRepository.Schema))
	require.NoError(t, err)
}

// NewTestContact creates a Contact domain object with basic test data (plaintext fields only)
func NewTestContact(t *testing.T) *client.Contact {
	t.Helper()
	return &client.Contact{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Lastname:  "DOE",
		Firstname: "John",
		Email:     "john.doe@example.com",
		Phone:     "0612345678",
		Position:  "position",
		CreatedAt: time.Now(),
	}
}

// NewTestContact creates a Contact domain object with basic test data (plaintext fields only)
func NewTestContactEncx(t *testing.T) *client.ContactEncx {
	t.Helper()
	return &client.ContactEncx{
		ID:                 uuid.New(),
		CreatedAt:          time.Now(),
		ClientIDEncrypted:  []byte("client_id_encrypted"),
		ClientIDHash:       "client_id_hash",
		LastnameEncrypted:  []byte("lastname_encrypted"),
		FirstnameEncrypted: []byte("lastname_encrypted"),
		EmailEncrypted:     []byte("email_encrypted"),
		EmailHash:          "email_hash",
		PhoneEncrypted:     []byte("phone_encrypted"),
		PositionEncrypted:  []byte("position_encrypted"),
		DEKEncrypted:       []byte("dek_encrypted"),
		KeyVersion:         1,
	}
}

// InsertContactEncx creates a Contact domain object with basic test data (plaintext fields only)
func InsertContactEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, contactEncx client.ContactEncx) error {
	t.Helper()

	query := fmt.Sprintf(`
		INSERT INTO %s.contacts (
			id, client_id_hash, client_id_encrypted, lastname_encrypted, firstname_encrypted, email_hash, email_encrypted, 
			phone_encrypted, position_encrypted, dek_encrypted, key_version
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, clientRepository.Schema)

	_, err := pool.Exec(ctx, query,
		contactEncx.ID, contactEncx.ClientIDHash, contactEncx.ClientIDEncrypted, contactEncx.LastnameEncrypted, contactEncx.FirstnameEncrypted, contactEncx.EmailHash,
		contactEncx.EmailEncrypted, contactEncx.PhoneEncrypted, contactEncx.PositionEncrypted, contactEncx.DEKEncrypted, contactEncx.KeyVersion,
	)

	return err
}

// GetContactEncx gets a ContactEncx using the ContactEncx ID.
func GetContactEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, contactID uuid.UUID) (*client.ContactEncx, error) {
	t.Helper()

	var contactEncx client.ContactEncx

	query := fmt.Sprintf(`
		SELECT
			id, client_id_encrypted, client_id_hash, firstname_encrypted, lastname_encrypted, email_encrypted,
			email_hash, phone_encrypted, position_encrypted, dek_encrypted, key_version
		FROM %s.contacts
		WHERE id = $1
	`, clientRepository.Schema)

	err := pool.QueryRow(ctx, query, contactID).Scan(
		&contactEncx.ID, &contactEncx.ClientIDEncrypted, &contactEncx.ClientIDHash,
		&contactEncx.FirstnameEncrypted, &contactEncx.LastnameEncrypted, &contactEncx.EmailEncrypted,
		&contactEncx.EmailHash, &contactEncx.PhoneEncrypted,
		&contactEncx.PositionEncrypted, &contactEncx.DEKEncrypted,
		&contactEncx.KeyVersion,
	)

	return &contactEncx, err
}

// GetContactEncxByEmailHash gets a ContactEncx using the email hash.
func GetContactEncxByEmailHash(t *testing.T, ctx context.Context, pool *pgxpool.Pool, emailHash string) (*client.ContactEncx, error) {
	t.Helper()

	var contactEncx client.ContactEncx

	query := fmt.Sprintf(`
		SELECT
			id, client_id_encrypted, client_id_hash, firstname_encrypted, lastname_encrypted, email_encrypted,
			email_hash, phone_encrypted, position_encrypted, dek_encrypted, key_version
		FROM %s.contacts
		WHERE email_hash = $1
	`, clientRepository.Schema)

	err := pool.QueryRow(ctx, query, emailHash).Scan(
		&contactEncx.ID, &contactEncx.ClientIDEncrypted, &contactEncx.ClientIDHash,
		&contactEncx.FirstnameEncrypted, &contactEncx.LastnameEncrypted, &contactEncx.EmailEncrypted,
		&contactEncx.EmailHash, &contactEncx.PhoneEncrypted,
		&contactEncx.PositionEncrypted, &contactEncx.DEKEncrypted,
		&contactEncx.KeyVersion,
	)

	return &contactEncx, err
}

// CountContactsByClientIDHash returns the number of contacts for a client ID hash
func CountContactsByClientIDHash(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientIDHash string) (int, error) {
	t.Helper()

	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s.contacts WHERE client_id_hash = $1`, clientRepository.Schema)
	err := pool.QueryRow(ctx, query, clientIDHash).Scan(&count)
	return count, err
}

// CreateTestClientWithContact creates a client by inserting an initial contact
// This represents a client "existing" in the system since clients are identified by contacts
func CreateTestClientWithContact(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientID uuid.UUID, clientIDHash string) error {
	t.Helper()

	initialContact := &client.ContactEncx{
		ID:                 uuid.New(),
		CreatedAt:          time.Now(),
		ClientIDEncrypted:  []byte("initial_client_id_encrypted"),
		ClientIDHash:       clientIDHash,
		LastnameEncrypted:  []byte("initial_lastname_encrypted"),
		FirstnameEncrypted: []byte("initial_firstname_encrypted"),
		EmailEncrypted:     []byte("initial_email_encrypted"),
		EmailHash:          "initial_email_hash_" + clientID.String(),
		PhoneEncrypted:     []byte("initial_phone_encrypted"),
		PositionEncrypted:  []byte("initial_position_encrypted"),
		DEKEncrypted:       []byte("initial_dek_encrypted"),
		KeyVersion:         1,
	}

	return InsertContactEncx(t, ctx, pool, *initialContact)
}
