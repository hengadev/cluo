package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM estimates WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count estimates: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			   line_items, estimated_total, notes, accepted, accepted_at, accepted_by,
			   created_at, updated_at
		FROM estimates
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query estimates: %w", err)
	}
	defer rows.Close()

	var estimates []*document.Estimate
	for rows.Next() {
		var estimate document.Estimate
		var lineItemsJSON []byte

		err := rows.Scan(
			&estimate.ID, &estimate.CaseID, &estimate.ClientID, &estimate.Status,
			&estimate.EstimateNumber, &estimate.IssueDate, &estimate.ValidUntil,
			&lineItemsJSON, &estimate.EstimatedTotal, &estimate.Notes,
			&estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
			&estimate.CreatedAt, &estimate.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan estimate: %w", err)
		}

		if err := json.Unmarshal(lineItemsJSON, &estimate.LineItems); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal line items: %w", err)
		}

		estimates = append(estimates, &estimate)
	}

	return estimates, total, nil
}

