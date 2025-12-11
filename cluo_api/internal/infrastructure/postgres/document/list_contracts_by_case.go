package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ListContractsByCase retrieves all contracts for a specific case with pagination.
func (r *Repository) ListContractsByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.Contract, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM contracts WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, contract_number, start_date, end_date,
			   scope_of_services, payment_terms, confidentiality, termination_clause,
			   signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law,
			   created_at, updated_at
		FROM contracts
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query contracts: %w", err)
	}
	defer rows.Close()

	var contracts []*document.Contract
	for rows.Next() {
		var contract document.Contract
		var signaturesJSON []byte

		err := rows.Scan(
			&contract.ID, &contract.CaseID, &contract.ClientID, &contract.Status,
			&contract.ContractNumber, &contract.StartDate, &contract.EndDate,
			&contract.ScopeOfServices, &contract.PaymentTerms, &contract.Confidentiality,
			&contract.TerminationClause, &signaturesJSON, &contract.LinkedMandateID,
			&contract.ContractValue, &contract.Currency, &contract.RenewalTerms, &contract.GoverningLaw,
			&contract.CreatedAt, &contract.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan contract: %w", err)
		}

		if err := json.Unmarshal(signaturesJSON, &contract.Signatures); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal signatures: %w", err)
		}

		contracts = append(contracts, &contract)
	}

	return contracts, total, nil
}
