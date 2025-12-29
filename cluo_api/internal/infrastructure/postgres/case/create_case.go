package caseRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/case"
)

func (r *Repository) CreateCase(ctx context.Context, caseEncx *caseDomain.CaseEncx) error {
	if caseEncx == nil {
		return fmt.Errorf("case cannot be nil")
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.cases (
			id, client_id, assigned_contact_id, case_subject_id, case_type, created_at,
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
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.CaseSubjectID,
		caseEncx.CaseType,
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
	if err != nil {
		return errs.ClassifyPgError("create case", err)
	}

	return nil
}

