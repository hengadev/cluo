package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateContract updates an existing contract in the database.
func (r *Repository) UpdateContract(ctx context.Context, contract *document.Contract) error {
	signaturesJSON, err := json.Marshal(contract.Signatures)
	if err != nil {
		return fmt.Errorf("failed to marshal signatures: %w", err)
	}

	query := `
		UPDATE contracts SET
			case_id = $2, client_id = $3, status = $4, contract_number = $5,
			start_date = $6, end_date = $7, scope_of_services = $8, payment_terms = $9,
			confidentiality = $10, termination_clause = $11, signatures = $12,
			linked_mandate_id = $13, contract_value = $14, currency = $15,
			renewal_terms = $16, governing_law = $17, updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, query,
		contract.ID, contract.CaseID, contract.ClientID, contract.Status,
		contract.ContractNumber, contract.StartDate, contract.EndDate,
		contract.ScopeOfServices, contract.PaymentTerms, contract.Confidentiality,
		contract.TerminationClause, signaturesJSON, contract.LinkedMandateID,
		contract.ContractValue, contract.Currency, contract.RenewalTerms, contract.GoverningLaw,
	)

	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	return nil
}
