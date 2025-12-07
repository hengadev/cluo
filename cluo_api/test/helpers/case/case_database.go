package caseHelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	caseRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/case"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearCasesTable truncates the cases table for clean test state
func ClearCasesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.cases RESTART IDENTITY CASCADE", caseRepository.Schema))
	require.NoError(t, err)
}

// InsertCaseEncx inserts a CaseEncx record into the database for testing
func InsertCaseEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseEncx *caseDomain.CaseEncx) error {
	t.Helper()

	query := fmt.Sprintf(`
		INSERT INTO %s.cases (
			id, client_id, assigned_contact_id, created_at,
			title_encrypted, description_encrypted, status_encrypted,
			updated_at_encrypted, dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, caseRepository.Schema)

	_, err := pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.CreatedAt,
		caseEncx.TitleEncrypted,
		caseEncx.DescriptionEncrypted,
		caseEncx.StatusEncrypted,
		caseEncx.UpdatedAtEncrypted,
		caseEncx.DEKEncrypted,
		caseEncx.KeyVersion,
		caseEncx.Metadata,
	)

	return err
}

// GetCaseEncxByID retrieves a case by ID from the database for testing
func GetCaseEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID uuid.UUID) (*caseDomain.CaseEncx, error) {
	t.Helper()

	query := fmt.Sprintf(`
		SELECT
			id, client_id, assigned_contact_id, created_at,
			title_encrypted, description_encrypted, status_encrypted,
			updated_at_encrypted, dek_encrypted, key_version, metadata
		FROM %s.cases WHERE id = $1
	`, caseRepository.Schema)

	caseEncx := &caseDomain.CaseEncx{}

	err := pool.QueryRow(ctx, query, caseID).Scan(
		&caseEncx.ID, &caseEncx.ClientID, &caseEncx.AssignedContactID, &caseEncx.CreatedAt,
		&caseEncx.TitleEncrypted, &caseEncx.DescriptionEncrypted, &caseEncx.StatusEncrypted,
		&caseEncx.UpdatedAtEncrypted, &caseEncx.DEKEncrypted, &caseEncx.KeyVersion, &caseEncx.Metadata,
	)

	return caseEncx, err
}

// CountCasesByClientID returns the number of cases for a client ID
func CountCasesByClientID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientID string) (int, error) {
	t.Helper()

	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s.cases WHERE client_id = $1`, caseRepository.Schema)
	err := pool.QueryRow(ctx, query, clientID).Scan(&count)
	return count, err
}

// CreateTestCaseWithClientID creates a case using a client ID that supposedly is in the database.
func CreateTestCaseWithClientID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientID string) error {
	t.Helper()

	initialCase := &caseDomain.CaseEncx{
		ID:                   uuid.New(),
		CreatedAt:            time.Now(),
		ClientID:             clientID,
		AssignedContactID:    func() *string { s := uuid.New().String(); return &s }(),
		TitleEncrypted:       []byte("initial_title_encrypted"),
		DescriptionEncrypted: []byte("initial_description_encrypted"),
		StatusEncrypted:      []byte("initial_status_encrypted"),
		UpdatedAtEncrypted:   []byte("initial_updatedat_encrypted"),
		DEKEncrypted:         []byte("initial_dek_encrypted"),
		KeyVersion:           1,
	}

	return InsertCaseEncx(t, ctx, pool, initialCase)
}
