package documentRepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetEstimateByID retrieves an estimate by its ID.
func (r *Repository) GetEstimateByID(ctx context.Context, id string) (*document.EstimateEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   estimatenumber_encrypted, lineitems_encrypted, estimatedtotal_encrypted, notes_encrypted,
			   issue_date, valid_until, accepted, accepted_at, accepted_by,
			   dek_encrypted, key_version, metadata
		FROM estimates
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var estimate document.EstimateEncx

	err := row.Scan(
		&estimate.ID, &estimate.Status, &estimate.CreatedAt, &estimate.UpdatedAt,
		&estimate.CaseIDEncrypted, &estimate.ClientIDEncrypted,
		&estimate.EstimateNumberEncrypted, &estimate.LineItemsEncrypted, &estimate.EstimatedTotalEncrypted, &estimate.NotesEncrypted,
		&estimate.IssueDate, &estimate.ValidUntil, &estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
		&estimate.DEKEncrypted, &estimate.KeyVersion, &estimate.Metadata,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("estimate not found")
		}
		return nil, fmt.Errorf("failed to get estimate: %w", err)
	}

	return &estimate, nil
}
