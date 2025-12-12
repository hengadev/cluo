package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateContract updates an existing contract in the database.
func (r *Repository) UpdateContract(ctx context.Context, contract *document.ContractEncx) error {
	query := `
		UPDATE contracts SET
			status = $2, updated_at = $3,
			caseid_encrypted = $4, clientid_encrypted = $5,
			contractnumber_encrypted = $6, scopeofservices_encrypted = $7, paymentterms_encrypted = $8,
			confidentiality_encrypted = $9, terminationclause_encrypted = $10, signatures_encrypted = $11,
			contractvalue_encrypted = $12, renewalterms_encrypted = $13,
			start_date = $14, end_date = $15, linked_mandate_id = $16, currency = $17, governing_law = $18,
			dek_encrypted = $19, key_version = $20, metadata = $21
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		contract.ID, contract.Status, contract.UpdatedAt,
		contract.CaseIDEncrypted, contract.ClientIDEncrypted,
		contract.ContractNumberEncrypted, contract.ScopeOfServicesEncrypted, contract.PaymentTermsEncrypted,
		contract.ConfidentialityEncrypted, contract.TerminationClauseEncrypted, contract.SignaturesEncrypted,
		contract.ContractValueEncrypted, contract.RenewalTermsEncrypted,
		contract.StartDate, contract.EndDate, contract.LinkedMandateID, contract.Currency, contract.GoverningLaw,
		contract.DEKEncrypted, contract.KeyVersion, contract.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	return nil
}
