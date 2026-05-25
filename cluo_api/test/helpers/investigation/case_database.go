package investigationHelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/investigation"
	investigationRepository "github.com/hengadev/cluo_api/internal/infrastructure/postgres/investigation"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearCasesTable truncates the cases table for clean test state
func ClearCasesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s.cases RESTART IDENTITY CASCADE", investigationRepository.Schema))
	require.NoError(t, err)
}

// InsertCaseEncx inserts a CaseEncx record into the database for testing
func InsertCaseEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseEncx *investigation.InvestigationEncx) error {
	t.Helper()

	query := fmt.Sprintf(`
		INSERT INTO %s.cases (
			id, client_id, assigned_contact_id, case_subject_id, case_type_id, created_at,
			title_encrypted, description_encrypted, external_reference_encrypted, external_reference_hash, status_encrypted,
			placename_encrypted, placename_hash,
			address1_encrypted, address1_hash,
			address2_encrypted, address2_hash,
			city_encrypted, city_hash,
			postal_code_encrypted, postal_code_hash,
			country_encrypted, country_hash,
			latitude_encrypted, latitude_hash,
			longitude_encrypted, longitude_hash,
			location_type_encrypted, location_type_hash,
			location_notes_encrypted, location_notes_hash,
			updated_at_encrypted, dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35)
	`, investigationRepository.Schema)

	_, err := pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.CaseSubjectID,
		caseEncx.CaseTypeID,
		caseEncx.CreatedAt,
		caseEncx.TitleEncrypted,
		caseEncx.DescriptionEncrypted,
		caseEncx.ExternalReferenceEncrypted,
		caseEncx.ExternalReferenceHash,
		caseEncx.StatusEncrypted,
		caseEncx.PlacenameEncrypted,
		caseEncx.PlacenameHash,
		caseEncx.Address1Encrypted,
		caseEncx.Address1Hash,
		caseEncx.Address2Encrypted,
		caseEncx.Address2Hash,
		caseEncx.CityEncrypted,
		caseEncx.CityHash,
		caseEncx.PostalCodeEncrypted,
		caseEncx.PostalCodeHash,
		caseEncx.CountryEncrypted,
		caseEncx.CountryHash,
		caseEncx.LatitudeEncrypted,
		caseEncx.LatitudeHash,
		caseEncx.LongitudeEncrypted,
		caseEncx.LongitudeHash,
		caseEncx.LocationTypeEncrypted,
		caseEncx.LocationTypeHash,
		caseEncx.LocationNotesEncrypted,
		caseEncx.LocationNotesHash,
		caseEncx.UpdatedAtEncrypted,
		caseEncx.DEKEncrypted,
		caseEncx.KeyVersion,
		caseEncx.Metadata,
	)

	return err
}

// GetCaseEncxByID retrieves a case by ID from the database for testing
func GetCaseEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID uuid.UUID) (*investigation.InvestigationEncx, error) {
	t.Helper()

	query := fmt.Sprintf(`
		SELECT
			id, client_id, assigned_contact_id, case_subject_id, case_type_id, created_at,
			title_encrypted, description_encrypted, external_reference_encrypted, external_reference_hash, status_encrypted,
			placename_encrypted, placename_hash,
			address1_encrypted, address1_hash,
			address2_encrypted, address2_hash,
			city_encrypted, city_hash,
			postal_code_encrypted, postal_code_hash,
			country_encrypted, country_hash,
			latitude_encrypted, latitude_hash,
			longitude_encrypted, longitude_hash,
			location_type_encrypted, location_type_hash,
			location_notes_encrypted, location_notes_hash,
			updated_at_encrypted, dek_encrypted, key_version, metadata
		FROM %s.cases WHERE id = $1
	`, investigationRepository.Schema)

	caseEncx := &investigation.InvestigationEncx{}

	err := pool.QueryRow(ctx, query, caseID).Scan(
		&caseEncx.ID, &caseEncx.ClientID, &caseEncx.AssignedContactID, &caseEncx.CaseSubjectID, &caseEncx.CaseTypeID, &caseEncx.CreatedAt,
		&caseEncx.TitleEncrypted, &caseEncx.DescriptionEncrypted, &caseEncx.ExternalReferenceEncrypted, &caseEncx.ExternalReferenceHash, &caseEncx.StatusEncrypted,
		&caseEncx.PlacenameEncrypted, &caseEncx.PlacenameHash,
		&caseEncx.Address1Encrypted, &caseEncx.Address1Hash,
		&caseEncx.Address2Encrypted, &caseEncx.Address2Hash,
		&caseEncx.CityEncrypted, &caseEncx.CityHash,
		&caseEncx.PostalCodeEncrypted, &caseEncx.PostalCodeHash,
		&caseEncx.CountryEncrypted, &caseEncx.CountryHash,
		&caseEncx.LatitudeEncrypted, &caseEncx.LatitudeHash,
		&caseEncx.LongitudeEncrypted, &caseEncx.LongitudeHash,
		&caseEncx.LocationTypeEncrypted, &caseEncx.LocationTypeHash,
		&caseEncx.LocationNotesEncrypted, &caseEncx.LocationNotesHash,
		&caseEncx.UpdatedAtEncrypted, &caseEncx.DEKEncrypted, &caseEncx.KeyVersion, &caseEncx.Metadata,
	)

	return caseEncx, err
}

// CountCasesByClientID returns the number of cases for a client ID
func CountCasesByClientID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientID string) (int, error) {
	t.Helper()

	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s.cases WHERE client_id = $1`, investigationRepository.Schema)
	err := pool.QueryRow(ctx, query, clientID).Scan(&count)
	return count, err
}

// CreateTestCaseWithClientID creates a case using a client ID that supposedly is in the database.
func CreateTestCaseWithClientID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, clientID uuid.UUID) error {
	t.Helper()

	contactID := uuid.New()
	initialCase := &investigation.InvestigationEncx{
		ID:                         uuid.New(),
		CreatedAt:                  time.Now(),
		ClientID:                   clientID,
		AssignedContactID:          &contactID,
		CaseSubjectID:              nil,
		CaseTypeID:                 nil,
		TitleEncrypted:             []byte("initial_title_encrypted"),
		DescriptionEncrypted:       []byte("initial_description_encrypted"),
		ExternalReferenceEncrypted: []byte("initial_external_ref_encrypted"),
		ExternalReferenceHash:      "initial_external_ref_hash",
		StatusEncrypted:            []byte("initial_status_encrypted"),
		PlacenameEncrypted:         []byte("initial_placename_encrypted"),
		PlacenameHash:              "initial_placename_hash",
		Address1Encrypted:          []byte("initial_address1_encrypted"),
		Address1Hash:               "initial_address1_hash",
		Address2Encrypted:          []byte("initial_address2_encrypted"),
		Address2Hash:               "initial_address2_hash",
		CityEncrypted:              []byte("initial_city_encrypted"),
		CityHash:                   "initial_city_hash",
		PostalCodeEncrypted:        []byte("initial_postal_code_encrypted"),
		PostalCodeHash:             "initial_postal_code_hash",
		CountryEncrypted:           []byte("initial_country_encrypted"),
		CountryHash:                "initial_country_hash",
		LatitudeEncrypted:          []byte("initial_latitude_encrypted"),
		LatitudeHash:               "initial_latitude_hash",
		LongitudeEncrypted:         []byte("initial_longitude_encrypted"),
		LongitudeHash:              "initial_longitude_hash",
		LocationTypeEncrypted:      []byte("initial_location_type_encrypted"),
		LocationTypeHash:           "initial_location_type_hash",
		LocationNotesEncrypted:     []byte("initial_location_notes_encrypted"),
		LocationNotesHash:          "initial_location_notes_hash",
		UpdatedAtEncrypted:         []byte("initial_updatedat_encrypted"),
		DEKEncrypted:               []byte("initial_dek_encrypted"),
		KeyVersion:                 1,
	}

	return InsertCaseEncx(t, ctx, pool, initialCase)
}
