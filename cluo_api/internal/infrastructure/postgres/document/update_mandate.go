package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateMandate updates an existing mandate in the database.
func (r *Repository) UpdateMandate(ctx context.Context, mandate *document.Mandate) error {
	var clientSignatureJSON, investigatorSignatureJSON []byte
	var err error

	if mandate.ClientSignature != nil {
		clientSignatureJSON, err = json.Marshal(mandate.ClientSignature)
		if err != nil {
			return fmt.Errorf("failed to marshal client signature: %w", err)
		}
	}

	if mandate.InvestigatorSignature != nil {
		investigatorSignatureJSON, err = json.Marshal(mandate.InvestigatorSignature)
		if err != nil {
			return fmt.Errorf("failed to marshal investigator signature: %w", err)
		}
	}

	query := `
		UPDATE mandates SET
			case_id = $2, client_id = $3, status = $4, mandate_number = $5,
			issue_date = $6, scope_of_work = $7, valid_from = $8, valid_until = $9,
			terms_conditions = $10, client_signature = $11, investigator_signature = $12,
			linked_estimate_id = $13, special_instructions = $14, jurisdiction = $15,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, query,
		mandate.ID, mandate.CaseID, mandate.ClientID, mandate.Status,
		mandate.MandateNumber, mandate.IssueDate, mandate.ScopeOfWork,
		mandate.ValidFrom, mandate.ValidUntil, mandate.TermsConditions,
		clientSignatureJSON, investigatorSignatureJSON, mandate.LinkedEstimateID,
		mandate.SpecialInstructions, mandate.Jurisdiction,
	)

	if err != nil {
		return fmt.Errorf("failed to update mandate: %w", err)
	}

	return nil
}
