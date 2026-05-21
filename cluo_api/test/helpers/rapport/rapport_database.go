package rapportHelpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const schema = "cases"

// ClearRapportsTable truncates the rapports table for clean test state.
func ClearRapportsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.rapports RESTART IDENTITY CASCADE", schema))
	require.NoError(t, err)
}

// InsertRapportEncx inserts a RapportEncx record into the database for testing.
func InsertRapportEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, rEncx *rapport.RapportEncx) error {
	t.Helper()

	query := fmt.Sprintf(`
		INSERT INTO %s.rapports (
			id, case_id, content_encrypted, dek_encrypted, key_version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, schema)

	_, err := pool.Exec(ctx, query,
		rEncx.ID,
		rEncx.CaseID,
		rEncx.ContentEncrypted,
		rEncx.DEKEncrypted,
		rEncx.KeyVersion,
		rEncx.CreatedAt,
		rEncx.UpdatedAt,
	)
	return err
}

// GetRapportEncxByCaseID retrieves a rapport by case ID from the database for testing.
// Returns nil, error when not found.
func GetRapportEncxByCaseID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID uuid.UUID) (*rapport.RapportEncx, error) {
	t.Helper()

	query := fmt.Sprintf(`
		SELECT id, case_id, content_encrypted, dek_encrypted, key_version, created_at, updated_at
		FROM %s.rapports
		WHERE case_id = $1
	`, schema)

	rEncx := &rapport.RapportEncx{}
	err := pool.QueryRow(ctx, query, caseID).Scan(
		&rEncx.ID,
		&rEncx.CaseID,
		&rEncx.ContentEncrypted,
		&rEncx.DEKEncrypted,
		&rEncx.KeyVersion,
		&rEncx.CreatedAt,
		&rEncx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rEncx, nil
}
