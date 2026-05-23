package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateEstimate creates a new estimate in the database.
func (r *Repository) CreateEstimate(ctx context.Context, estimate *document.EstimateEncx) error {
	query := `
		INSERT INTO estimates (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted,
			estimatenumber_encrypted, lineitems_encrypted, estimatedtotal_encrypted, notes_encrypted,
			issue_date, valid_until, accepted, accepted_at, accepted_by,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.pool.Exec(ctx, query,
		estimate.ID, estimate.Status, estimate.CreatedAt, estimate.UpdatedAt,
		estimate.CaseIDEncrypted, estimate.ClientIDEncrypted,
		estimate.EstimateNumberEncrypted, estimate.LineItemsEncrypted, estimate.EstimatedTotalEncrypted, estimate.NotesEncrypted,
		estimate.IssueDate, estimate.ValidUntil, estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
		estimate.DEKEncrypted, estimate.KeyVersion, estimate.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to create estimate: %w", err)
	}

	return nil
}
