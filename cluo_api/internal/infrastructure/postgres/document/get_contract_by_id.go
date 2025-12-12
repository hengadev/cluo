package documentRepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetContractByID retrieves a contract by its ID.
func (r *Repository) GetContractByID(ctx context.Context, id string) (*document.ContractEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   contractnumber_encrypted, scopeofservices_encrypted, paymentterms_encrypted,
			   confidentiality_encrypted, terminationclause_encrypted, signatures_encrypted,
			   contractvalue_encrypted, renewalterms_encrypted,
			   start_date, end_date, linked_mandate_id, currency, governing_law,
			   dek_encrypted, key_version, metadata
		FROM contracts
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var contract document.ContractEncx

	err := row.Scan(
		&contract.ID, &contract.Status, &contract.CreatedAt, &contract.UpdatedAt,
		&contract.CaseIDEncrypted, &contract.ClientIDEncrypted,
		&contract.ContractNumberEncrypted, &contract.ScopeOfServicesEncrypted, &contract.PaymentTermsEncrypted,
		&contract.ConfidentialityEncrypted, &contract.TerminationClauseEncrypted, &contract.SignaturesEncrypted,
		&contract.ContractValueEncrypted, &contract.RenewalTermsEncrypted,
		&contract.StartDate, &contract.EndDate, &contract.LinkedMandateID, &contract.Currency, &contract.GoverningLaw,
		&contract.DEKEncrypted, &contract.KeyVersion, &contract.Metadata,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract not found")
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	return &contract, nil
}
