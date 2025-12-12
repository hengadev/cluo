package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ListContractsByCase retrieves all contracts for a specific case with pagination.
func (r *Repository) ListContractsByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.ContractEncx, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM contracts WHERE caseid_encrypted = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   contractnumber_encrypted, scopeofservices_encrypted, paymentterms_encrypted,
			   confidentiality_encrypted, terminationclause_encrypted, signatures_encrypted,
			   contractvalue_encrypted, renewalterms_encrypted,
			   start_date, end_date, linked_mandate_id, currency, governing_law,
			   dek_encrypted, key_version, metadata
		FROM contracts
		WHERE caseid_encrypted = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query contracts: %w", err)
	}
	defer rows.Close()

	var contracts []*document.ContractEncx
	for rows.Next() {
		var contract document.ContractEncx

		err := rows.Scan(
			&contract.ID, &contract.Status, &contract.CreatedAt, &contract.UpdatedAt,
			&contract.CaseIDEncrypted, &contract.ClientIDEncrypted,
			&contract.ContractNumberEncrypted, &contract.ScopeOfServicesEncrypted, &contract.PaymentTermsEncrypted,
			&contract.ConfidentialityEncrypted, &contract.TerminationClauseEncrypted, &contract.SignaturesEncrypted,
			&contract.ContractValueEncrypted, &contract.RenewalTermsEncrypted,
			&contract.StartDate, &contract.EndDate, &contract.LinkedMandateID, &contract.Currency, &contract.GoverningLaw,
			&contract.DEKEncrypted, &contract.KeyVersion, &contract.Metadata,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan contract: %w", err)
		}

		contracts = append(contracts, &contract)
	}

	return contracts, total, nil
}
