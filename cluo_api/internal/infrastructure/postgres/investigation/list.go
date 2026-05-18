package investigationRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

// List implements the general List function with filtering
func (r *Repository) List(ctx context.Context, f investigation.Filter, p investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	// Validate pagination
	if err := p.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	// Build WHERE clauses for fields that can be filtered at database level
	whereClauses := []string{"1=1"} // Start with true condition
	args := []interface{}{}
	argIndex := 1

	// Add filter conditions for non-encrypted fields
	if f.ClientID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("client_id = $%d", argIndex))
		args = append(args, *f.ClientID)
		argIndex++
	}

	if f.AssignedContactID != nil {
		if *f.AssignedContactID == (uuid.UUID{}) { // Check for nil UUID
			whereClauses = append(whereClauses, "assigned_contact_id IS NULL")
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("assigned_contact_id = $%d", argIndex))
			args = append(args, *f.AssignedContactID)
			argIndex++
		}
	}

	if f.CaseSubjectID != nil {
		if *f.CaseSubjectID == (uuid.UUID{}) { // Check for nil UUID
			whereClauses = append(whereClauses, "case_subject_id IS NULL")
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("case_subject_id = $%d", argIndex))
			args = append(args, *f.CaseSubjectID)
			argIndex++
		}
	}

	if f.CaseType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("case_type = $%d", argIndex))
		args = append(args, *f.CaseType)
		argIndex++
	}

	// Add location hash filters (hashes computed by application layer)
	if f.CityHash != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("city_hash = $%d", argIndex))
		args = append(args, *f.CityHash)
		argIndex++
	}

	if f.PostalCodeHash != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("postal_code_hash = $%d", argIndex))
		args = append(args, *f.PostalCodeHash)
		argIndex++
	}

	if f.CountryHash != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("country_hash = $%d", argIndex))
		args = append(args, *f.CountryHash)
		argIndex++
	}

	if f.DateCreatedFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *f.DateCreatedFrom)
		argIndex++
	}

	if f.DateCreatedTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *f.DateCreatedTo)
		argIndex++
	}

	// Note: Status, Search, and UpdatedAt filtering require decryption and should be handled at service layer

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.cases
		%s
	`, r.schema, whereSQL)

	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cases: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT
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
		FROM %s.cases
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, r.schema, whereSQL, argIndex, argIndex+1)

	args = append(args, p.PageSize, p.GetOffset())

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query cases: %w", err)
	}
	defer rows.Close()

	var cases []*investigation.InvestigationEncx
	for rows.Next() {
		caseEncx := &investigation.InvestigationEncx{}
		err := rows.Scan(
			&caseEncx.ID, &caseEncx.ClientID, &caseEncx.AssignedContactID, &caseEncx.CaseSubjectID, &caseEncx.CaseType, &caseEncx.CreatedAt,
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
			return nil, 0, fmt.Errorf("failed to scan case row: %w", err)
		}

		cases = append(cases, caseEncx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating case rows: %w", err)
	}

	return cases, total, nil
}

