package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	signaturesJSON, err := json.Marshal(contract.Signatures)
	if err != nil {
		return fmt.Errorf("failed to marshal signatures: %w", err)
	}

	query := `
		INSERT INTO contracts (
			id, case_id, client_id, status, contract_number, start_date, end_date,
			scope_of_services, payment_terms, confidentiality, termination_clause,
			signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err = r.pool.Exec(ctx, query,
		contract.ID, contract.CaseID, contract.ClientID, contract.Status,
		contract.ContractNumber, contract.StartDate, contract.EndDate,
		contract.ScopeOfServices, contract.PaymentTerms, contract.Confidentiality,
		contract.TerminationClause, signaturesJSON, contract.LinkedMandateID,
		contract.ContractValue, contract.Currency, contract.RenewalTerms, contract.GoverningLaw,
	)

	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	return nil
}
