package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateEstimate updates an existing estimate in the database.
func (r *Repository) UpdateEstimate(ctx context.Context, estimate *document.Estimate) error {
	lineItemsJSON, err := json.Marshal(estimate.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		UPDATE estimates SET
			case_id = $2, client_id = $3, status = $4, estimate_number = $5,
			issue_date = $6, valid_until = $7, line_items = $8, estimated_total = $9,
			notes = $10, accepted = $11, accepted_at = $12, accepted_by = $13,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, query,
		estimate.ID, estimate.CaseID, estimate.ClientID, estimate.Status,
		estimate.EstimateNumber, estimate.IssueDate, estimate.ValidUntil,
		lineItemsJSON, estimate.EstimatedTotal, estimate.Notes,
		estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update estimate: %w", err)
	}

	return nil
}
