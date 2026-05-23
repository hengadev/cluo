package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateMandate updates an existing mandate in the database.
func (r *Repository) UpdateMandate(ctx context.Context, mandate *document.MandateEncx) error {
	query := `
		UPDATE mandates SET
			status = $2, updated_at = $3,
			caseid_encrypted = $4, clientid_encrypted = $5,
			mandatenumber_encrypted = $6, scopeofwork_encrypted = $7, termsconditions_encrypted = $8,
			clientsignature_encrypted = $9, investigatorsignature_encrypted = $10, specialinstructions_encrypted = $11,
			issue_date = $12, valid_from = $13, valid_until = $14, linked_estimate_id = $15, jurisdiction = $16,
			dek_encrypted = $17, key_version = $18, metadata = $19
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		mandate.ID, mandate.Status, mandate.UpdatedAt,
		mandate.CaseIDEncrypted, mandate.ClientIDEncrypted,
		mandate.MandateNumberEncrypted, mandate.ScopeOfWorkEncrypted, mandate.TermsConditionsEncrypted,
		mandate.ClientSignatureEncrypted, mandate.InvestigatorSignatureEncrypted, mandate.SpecialInstructionsEncrypted,
		mandate.IssueDate, mandate.ValidFrom, mandate.ValidUntil, mandate.LinkedEstimateID, mandate.Jurisdiction,
		mandate.DEKEncrypted, mandate.KeyVersion, mandate.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to update mandate: %w", err)
	}

	return nil
}
