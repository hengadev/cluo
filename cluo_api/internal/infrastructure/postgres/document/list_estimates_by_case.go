package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ListEstimatesByCase retrieves all estimates for a specific case with pagination.
func (r *Repository) ListEstimatesByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.EstimateEncx, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM estimates WHERE caseid_encrypted = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count estimates: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   estimatenumber_encrypted, lineitems_encrypted, estimatedtotal_encrypted, notes_encrypted,
			   issue_date, valid_until, accepted, accepted_at, accepted_by,
			   dek_encrypted, key_version, metadata
		FROM estimates
		WHERE caseid_encrypted = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query estimates: %w", err)
	}
	defer rows.Close()

	var estimates []*document.EstimateEncx
	for rows.Next() {
		var estimate document.EstimateEncx

		err := rows.Scan(
			&estimate.ID, &estimate.Status, &estimate.CreatedAt, &estimate.UpdatedAt,
			&estimate.CaseIDEncrypted, &estimate.ClientIDEncrypted,
			&estimate.EstimateNumberEncrypted, &estimate.LineItemsEncrypted, &estimate.EstimatedTotalEncrypted, &estimate.NotesEncrypted,
			&estimate.IssueDate, &estimate.ValidUntil, &estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
			&estimate.DEKEncrypted, &estimate.KeyVersion, &estimate.Metadata,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan estimate: %w", err)
		}

		estimates = append(estimates, &estimate)
	}

	return estimates, total, nil
}
