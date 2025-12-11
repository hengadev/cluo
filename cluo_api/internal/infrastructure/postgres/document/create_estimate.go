package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateEstimate creates a new estimate in the database.
func (r *Repository) CreateEstimate(ctx context.Context, estimate *document.Estimate) error {
	lineItemsJSON, err := json.Marshal(estimate.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		INSERT INTO estimates (
			id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			line_items, estimated_total, notes, accepted, accepted_at, accepted_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = r.pool.Exec(ctx, query,
		estimate.ID, estimate.CaseID, estimate.ClientID, estimate.Status,
		estimate.EstimateNumber, estimate.IssueDate, estimate.ValidUntil,
		lineItemsJSON, estimate.EstimatedTotal, estimate.Notes,
		estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create estimate: %w", err)
	}

	return nil
}
