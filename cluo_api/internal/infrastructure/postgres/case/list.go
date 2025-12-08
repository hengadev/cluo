package caseRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

// List implements the general List function with filtering
func (r *Repository) List(ctx context.Context, f caseDomain.CaseFilter, p caseDomain.Pagination) ([]*caseDomain.CaseEncx, int, error) {
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
			id,
			client_id,
			assigned_contact_id,
			created_at,
			title_encrypted,
			description_encrypted,
			status_encrypted,
			updated_at_encrypted,
			dek_encrypted,
			key_version,
			metadata
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

	var cases []*caseDomain.CaseEncx
	for rows.Next() {
		caseEncx := &caseDomain.CaseEncx{}
		err := rows.Scan(
			&caseEncx.ID,
			&caseEncx.ClientID,
			&caseEncx.AssignedContactID,
			&caseEncx.CreatedAt,
			&caseEncx.TitleEncrypted,
			&caseEncx.DescriptionEncrypted,
			&caseEncx.StatusEncrypted,
			&caseEncx.UpdatedAtEncrypted,
			&caseEncx.DEKEncrypted,
			&caseEncx.KeyVersion,
			&caseEncx.Metadata,
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

