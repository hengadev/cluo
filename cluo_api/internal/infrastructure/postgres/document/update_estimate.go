package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateEstimate updates an existing estimate in the database.
func (r *Repository) UpdateEstimate(ctx context.Context, estimate *document.EstimateEncx) error {
	query := `
		UPDATE estimates SET
			status = $2, updated_at = $3,
			caseid_encrypted = $4, clientid_encrypted = $5,
			estimatenumber_encrypted = $6, lineitems_encrypted = $7, estimatedtotal_encrypted = $8, notes_encrypted = $9,
			issue_date = $10, valid_until = $11, accepted = $12, accepted_at = $13, accepted_by = $14,
			dek_encrypted = $15, key_version = $16, metadata = $17
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		estimate.ID, estimate.Status, estimate.UpdatedAt,
		estimate.CaseIDEncrypted, estimate.ClientIDEncrypted,
		estimate.EstimateNumberEncrypted, estimate.LineItemsEncrypted, estimate.EstimatedTotalEncrypted, estimate.NotesEncrypted,
		estimate.IssueDate, estimate.ValidUntil, estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
		estimate.DEKEncrypted, estimate.KeyVersion, estimate.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to update estimate: %w", err)
	}

	return nil
}
