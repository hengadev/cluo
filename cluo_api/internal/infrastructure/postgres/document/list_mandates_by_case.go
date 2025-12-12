package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ListMandatesByCase retrieves all mandates for a specific case with pagination.
func (r *Repository) ListMandatesByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.MandateEncx, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM mandates WHERE caseid_encrypted = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mandates: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   mandatenumber_encrypted, scopeofwork_encrypted, termsconditions_encrypted,
			   clientsignature_encrypted, investigatorsignature_encrypted, specialinstructions_encrypted,
			   issue_date, valid_from, valid_until, linked_estimate_id, jurisdiction,
			   dek_encrypted, key_version, metadata
		FROM mandates
		WHERE caseid_encrypted = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query mandates: %w", err)
	}
	defer rows.Close()

	var mandates []*document.MandateEncx
	for rows.Next() {
		var mandate document.MandateEncx

		err := rows.Scan(
			&mandate.ID, &mandate.Status, &mandate.CreatedAt, &mandate.UpdatedAt,
			&mandate.CaseIDEncrypted, &mandate.ClientIDEncrypted,
			&mandate.MandateNumberEncrypted, &mandate.ScopeOfWorkEncrypted, &mandate.TermsConditionsEncrypted,
			&mandate.ClientSignatureEncrypted, &mandate.InvestigatorSignatureEncrypted, &mandate.SpecialInstructionsEncrypted,
			&mandate.IssueDate, &mandate.ValidFrom, &mandate.ValidUntil, &mandate.LinkedEstimateID, &mandate.Jurisdiction,
			&mandate.DEKEncrypted, &mandate.KeyVersion, &mandate.Metadata,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mandate: %w", err)
		}

		mandates = append(mandates, &mandate)
	}

	return mandates, total, nil
}
