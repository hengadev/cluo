package caseRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (r *Repository) UpdateCase(ctx context.Context, caseEncx *caseDomain.CaseEncx) error {
	if caseEncx == nil {
		return fmt.Errorf("case cannot be nil")
	}

	query := fmt.Sprintf(`
		UPDATE %s.cases SET
			client_id = $2,
			assigned_contact_id = $3,
			case_subject_id = $4,
			case_type = $5,
			title_encrypted = $6,
			description_encrypted = $7,
			external_reference_encrypted = $8,
			external_reference_hash = $9,
			status_encrypted = $10,
			placename_encrypted = $11,
			placename_hash = $12,
			address1_encrypted = $13,
			address1_hash = $14,
			address2_encrypted = $15,
			address2_hash = $16,
			city_encrypted = $17,
			city_hash = $18,
			postal_code_encrypted = $19,
			postal_code_hash = $20,
			country_encrypted = $21,
			country_hash = $22,
			latitude_encrypted = $23,
			latitude_hash = $24,
			longitude_encrypted = $25,
			longitude_hash = $26,
			location_type_encrypted = $27,
			location_type_hash = $28,
			location_notes_encrypted = $29,
			location_notes_hash = $30,
			updated_at_encrypted = $31,
			dek_encrypted = $32,
			key_version = $33,
			metadata = $34
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.CaseSubjectID,
		caseEncx.CaseType,
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
	if err != nil {
		return errs.ClassifyPgError("update case", err)
	}

	// Check if any row was actually updated
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "case for update")
	}

	return nil
}

