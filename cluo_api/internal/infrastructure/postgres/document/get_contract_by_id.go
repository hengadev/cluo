package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	query := `
		SELECT id, case_id, client_id, status, contract_number, start_date, end_date,
			   scope_of_services, payment_terms, confidentiality, termination_clause,
			   signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law,
			   created_at, updated_at
		FROM contracts
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var contract document.Contract
	var signaturesJSON []byte

	err := row.Scan(
		&contract.ID, &contract.CaseID, &contract.ClientID, &contract.Status,
		&contract.ContractNumber, &contract.StartDate, &contract.EndDate,
		&contract.ScopeOfServices, &contract.PaymentTerms, &contract.Confidentiality,
		&contract.TerminationClause, &signaturesJSON, &contract.LinkedMandateID,
		&contract.ContractValue, &contract.Currency, &contract.RenewalTerms, &contract.GoverningLaw,
		&contract.CreatedAt, &contract.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract not found")
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	if err := json.Unmarshal(signaturesJSON, &contract.Signatures); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signatures: %w", err)
	}

	return &contract, nil
}
