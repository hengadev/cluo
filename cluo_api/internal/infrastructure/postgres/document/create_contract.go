package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateContract creates a new contract in the database.
func (r *Repository) CreateContract(ctx context.Context, contract *document.ContractEncx) error {
	query := `
		INSERT INTO contracts (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted,
			contractnumber_encrypted, scopeofservices_encrypted, paymentterms_encrypted,
			confidentiality_encrypted, terminationclause_encrypted, signatures_encrypted,
			contractvalue_encrypted, renewalterms_encrypted,
			start_date, end_date, linked_mandate_id, currency, governing_law,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	_, err := r.pool.Exec(ctx, query,
		contract.ID, contract.Status, contract.CreatedAt, contract.UpdatedAt,
		contract.CaseIDEncrypted, contract.ClientIDEncrypted,
		contract.ContractNumberEncrypted, contract.ScopeOfServicesEncrypted, contract.PaymentTermsEncrypted,
		contract.ConfidentialityEncrypted, contract.TerminationClauseEncrypted, contract.SignaturesEncrypted,
		contract.ContractValueEncrypted, contract.RenewalTermsEncrypted,
		contract.StartDate, contract.EndDate, contract.LinkedMandateID, contract.Currency, contract.GoverningLaw,
		contract.DEKEncrypted, contract.KeyVersion, contract.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	return nil
}
