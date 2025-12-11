package documentRepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetEstimateByID retrieves an estimate by its ID.
func (r *Repository) GetEstimateByID(ctx context.Context, id string) (*document.Estimate, error) {
	query := `
		SELECT id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			   line_items, estimated_total, notes, accepted, accepted_at, accepted_by,
			   created_at, updated_at
		FROM estimates
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var estimate document.Estimate
	var lineItemsJSON []byte

	err := row.Scan(
		&estimate.ID, &estimate.CaseID, &estimate.ClientID, &estimate.Status,
		&estimate.EstimateNumber, &estimate.IssueDate, &estimate.ValidUntil,
		&lineItemsJSON, &estimate.EstimatedTotal, &estimate.Notes,
		&estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
		&estimate.CreatedAt, &estimate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("estimate not found")
		}
		return nil, fmt.Errorf("failed to get estimate: %w", err)
	}

	if err := json.Unmarshal(lineItemsJSON, &estimate.LineItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal line items: %w", err)
	}

	return &estimate, nil
}
