package investigationRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

// ListByClient returns cases for a specific client with pagination
func (r *Repository) ListByClient(ctx context.Context, clientID uuid.UUID, p investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	// Validate pagination
	if err := p.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	// Validate clientID
	if clientID == uuid.Nil {
		return nil, 0, fmt.Errorf("client ID cannot be nil")
	}

	// Get total count for this client
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.cases
		WHERE client_id = $1
	`, r.schema)

	var total int
	err := r.pool.QueryRow(ctx, countQuery, clientID).Scan(&total)
	if err != nil {
		return nil, 0, errs.ClassifyPgError("count cases by client", err)
	}

	// Get paginated results for this client
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
		FROM %s.cases
		WHERE client_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, clientID, p.PageSize, p.GetOffset())
	if err != nil {
		return nil, 0, errs.ClassifyPgError("query cases by client", err)
	}
	defer rows.Close()

	var cases []*investigation.InvestigationEncx
	for rows.Next() {
		caseEncx := &investigation.InvestigationEncx{}
		err := rows.Scan(
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
		if err != nil {
			return nil, 0, errs.ClassifyPgError("scan case row", err)
		}

		cases = append(cases, caseEncx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errs.ClassifyPgError("iterate case rows", err)
	}

	return cases, total, nil
}
