package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateMandate creates a new mandate in the database.
func (r *Repository) CreateMandate(ctx context.Context, mandate *document.Mandate) error {
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
		INSERT INTO mandates (
			id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			valid_from, valid_until, terms_conditions, client_signature,
			investigator_signature, linked_estimate_id, special_instructions, jurisdiction
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = r.pool.Exec(ctx, query,
		mandate.ID, mandate.CaseID, mandate.ClientID, mandate.Status,
		mandate.MandateNumber, mandate.IssueDate, mandate.ScopeOfWork,
		mandate.ValidFrom, mandate.ValidUntil, mandate.TermsConditions,
		clientSignatureJSON, investigatorSignatureJSON, mandate.LinkedEstimateID,
		mandate.SpecialInstructions, mandate.Jurisdiction,
	)

	if err != nil {
		return fmt.Errorf("failed to create mandate: %w", err)
	}

	return nil
}
