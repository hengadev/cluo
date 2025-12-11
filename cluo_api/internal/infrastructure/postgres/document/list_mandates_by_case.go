package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM mandates WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mandates: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			   valid_from, valid_until, terms_conditions, client_signature,
			   investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
			   created_at, updated_at
		FROM mandates
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query mandates: %w", err)
	}
	defer rows.Close()

	var mandates []*document.Mandate
	for rows.Next() {
		var mandate document.Mandate
		var clientSignatureJSON, investigatorSignatureJSON []byte

		err := rows.Scan(
			&mandate.ID, &mandate.CaseID, &mandate.ClientID, &mandate.Status,
			&mandate.MandateNumber, &mandate.IssueDate, &mandate.ScopeOfWork,
			&mandate.ValidFrom, &mandate.ValidUntil, &mandate.TermsConditions,
			&clientSignatureJSON, &investigatorSignatureJSON, &mandate.LinkedEstimateID,
			&mandate.SpecialInstructions, &mandate.Jurisdiction,
			&mandate.CreatedAt, &mandate.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mandate: %w", err)
		}

		if len(clientSignatureJSON) > 0 {
			if err := json.Unmarshal(clientSignatureJSON, &mandate.ClientSignature); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal client signature: %w", err)
			}
		}

		if len(investigatorSignatureJSON) > 0 {
			if err := json.Unmarshal(investigatorSignatureJSON, &mandate.InvestigatorSignature); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal investigator signature: %w", err)
			}
		}

		mandates = append(mandates, &mandate)
	}

	return mandates, total, nil
}

