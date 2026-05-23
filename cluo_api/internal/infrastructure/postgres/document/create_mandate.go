package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateMandate creates a new mandate in the database.
func (r *Repository) CreateMandate(ctx context.Context, mandate *document.MandateEncx) error {
	query := `
		INSERT INTO mandates (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted,
			mandatenumber_encrypted, scopeofwork_encrypted, termsconditions_encrypted,
			clientsignature_encrypted, investigatorsignature_encrypted, specialinstructions_encrypted,
			issue_date, valid_from, valid_until, linked_estimate_id, jurisdiction,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := r.pool.Exec(ctx, query,
		mandate.ID, mandate.Status, mandate.CreatedAt, mandate.UpdatedAt,
		mandate.CaseIDEncrypted, mandate.ClientIDEncrypted,
		mandate.MandateNumberEncrypted, mandate.ScopeOfWorkEncrypted, mandate.TermsConditionsEncrypted,
		mandate.ClientSignatureEncrypted, mandate.InvestigatorSignatureEncrypted, mandate.SpecialInstructionsEncrypted,
		mandate.IssueDate, mandate.ValidFrom, mandate.ValidUntil, mandate.LinkedEstimateID, mandate.Jurisdiction,
		mandate.DEKEncrypted, mandate.KeyVersion, mandate.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to create mandate: %w", err)
	}

	return nil
}
