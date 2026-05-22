package subjectHelpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	subjectDomain "github.com/hengadev/cluo_api/internal/domain/subject"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const schema = "cases"

// ClearCaseSubjectsTable truncates the case_subjects table for clean test state.
func ClearCaseSubjectsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.case_subjects RESTART IDENTITY CASCADE", schema))
	require.NoError(t, err)
}

// InsertSubjectEncx inserts a SubjectEncx record into the database for testing.
func InsertSubjectEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, sEncx *subjectDomain.SubjectEncx) error {
	t.Helper()

	query := fmt.Sprintf(`
		INSERT INTO %s.case_subjects (
			id, created_at,
			lastname_encrypted, lastname_hash,
			firstname_encrypted, firstname_hash,
			email_encrypted, email_hash,
			phone_encrypted,
			city_encrypted, city_hash,
			postal_code_encrypted, postal_code_hash,
			address1_encrypted, address1_hash,
			address2_encrypted, address2_hash,
			occupation_encrypted, occupation_hash,
			notes_encrypted,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	`, schema)

	_, err := pool.Exec(ctx, query,
		sEncx.ID,
		sEncx.CreatedAt,
		sEncx.LastnameEncrypted,
		sEncx.LastnameHash,
		sEncx.FirstnameEncrypted,
		sEncx.FirstnameHash,
		sEncx.EmailEncrypted,
		sEncx.EmailHash,
		sEncx.PhoneEncrypted,
		sEncx.CityEncrypted,
		sEncx.CityHash,
		sEncx.PostalCodeEncrypted,
		sEncx.PostalCodeHash,
		sEncx.Address1Encrypted,
		sEncx.Address1Hash,
		sEncx.Address2Encrypted,
		sEncx.Address2Hash,
		sEncx.OccupationEncrypted,
		sEncx.OccupationHash,
		sEncx.NotesEncrypted,
		sEncx.DEKEncrypted,
		sEncx.KeyVersion,
		sEncx.Metadata,
	)
	return err
}

// GetSubjectEncxByID retrieves a subject by ID from the database for testing.
// Returns nil, error when not found.
func GetSubjectEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) (*subjectDomain.SubjectEncx, error) {
	t.Helper()

	query := fmt.Sprintf(`
		SELECT
			id, created_at,
			lastname_encrypted, lastname_hash,
			firstname_encrypted, firstname_hash,
			email_encrypted, email_hash,
			phone_encrypted,
			city_encrypted, city_hash,
			postal_code_encrypted, postal_code_hash,
			address1_encrypted, address1_hash,
			address2_encrypted, address2_hash,
			occupation_encrypted, occupation_hash,
			notes_encrypted,
			dek_encrypted, key_version, metadata
		FROM %s.case_subjects
		WHERE id = $1
	`, schema)

	sEncx := &subjectDomain.SubjectEncx{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&sEncx.ID,
		&sEncx.CreatedAt,
		&sEncx.LastnameEncrypted,
		&sEncx.LastnameHash,
		&sEncx.FirstnameEncrypted,
		&sEncx.FirstnameHash,
		&sEncx.EmailEncrypted,
		&sEncx.EmailHash,
		&sEncx.PhoneEncrypted,
		&sEncx.CityEncrypted,
		&sEncx.CityHash,
		&sEncx.PostalCodeEncrypted,
		&sEncx.PostalCodeHash,
		&sEncx.Address1Encrypted,
		&sEncx.Address1Hash,
		&sEncx.Address2Encrypted,
		&sEncx.Address2Hash,
		&sEncx.OccupationEncrypted,
		&sEncx.OccupationHash,
		&sEncx.NotesEncrypted,
		&sEncx.DEKEncrypted,
		&sEncx.KeyVersion,
		&sEncx.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return sEncx, nil
}

// CountCaseSubjects counts the total number of subjects in the database.
func CountCaseSubjects(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s.case_subjects`, schema)
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)
	return count
}
