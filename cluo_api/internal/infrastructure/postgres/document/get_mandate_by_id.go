package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	query := `
		SELECT id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			   valid_from, valid_until, terms_conditions, client_signature,
			   investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
			   created_at, updated_at
		FROM mandates
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var mandate document.Mandate
	var clientSignatureJSON, investigatorSignatureJSON []byte

	err := row.Scan(
		&mandate.ID, &mandate.CaseID, &mandate.ClientID, &mandate.Status,
		&mandate.MandateNumber, &mandate.IssueDate, &mandate.ScopeOfWork,
		&mandate.ValidFrom, &mandate.ValidUntil, &mandate.TermsConditions,
		&clientSignatureJSON, &investigatorSignatureJSON, &mandate.LinkedEstimateID,
		&mandate.SpecialInstructions, &mandate.Jurisdiction,
		&mandate.CreatedAt, &mandate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mandate not found")
		}
		return nil, fmt.Errorf("failed to get mandate: %w", err)
	}

	if len(clientSignatureJSON) > 0 {
		if err := json.Unmarshal(clientSignatureJSON, &mandate.ClientSignature); err != nil {
			return nil, fmt.Errorf("failed to unmarshal client signature: %w", err)
		}
	}

	if len(investigatorSignatureJSON) > 0 {
		if err := json.Unmarshal(investigatorSignatureJSON, &mandate.InvestigatorSignature); err != nil {
			return nil, fmt.Errorf("failed to unmarshal investigator signature: %w", err)
		}
	}

	return &mandate, nil
}
