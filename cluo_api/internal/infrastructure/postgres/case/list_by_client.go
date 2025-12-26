package caseRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

// ListByClient returns cases for a specific client with pagination
func (r *Repository) ListByClient(ctx context.Context, clientID uuid.UUID, p caseDomain.Pagination) ([]*caseDomain.CaseEncx, int, error) {
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
			id,
			client_id,
			assigned_contact_id,
			case_type,
			created_at,
			title_encrypted,
			description_encrypted,
			external_reference_encrypted,
			status_encrypted,
			updated_at_encrypted,
			dek_encrypted,
			key_version,
			metadata
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

	var cases []*caseDomain.CaseEncx
	for rows.Next() {
		caseEncx := &caseDomain.CaseEncx{}
		err := rows.Scan(
			&caseEncx.ID,
			&caseEncx.ClientID,
			&caseEncx.AssignedContactID,
			&caseEncx.CaseType,
			&caseEncx.CreatedAt,
			&caseEncx.TitleEncrypted,
			&caseEncx.DescriptionEncrypted,
			&caseEncx.ExternalReferenceEncrypted,
			&caseEncx.StatusEncrypted,
			&caseEncx.UpdatedAtEncrypted,
			&caseEncx.DEKEncrypted,
			&caseEncx.KeyVersion,
			&caseEncx.Metadata,
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
