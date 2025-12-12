package documentRepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetMandateByID retrieves a mandate by its ID.
func (r *Repository) GetMandateByID(ctx context.Context, id string) (*document.MandateEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   mandatenumber_encrypted, scopeofwork_encrypted, termsconditions_encrypted,
			   clientsignature_encrypted, investigatorsignature_encrypted, specialinstructions_encrypted,
			   issue_date, valid_from, valid_until, linked_estimate_id, jurisdiction,
			   dek_encrypted, key_version, metadata
		FROM mandates
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var mandate document.MandateEncx

	err := row.Scan(
		&mandate.ID, &mandate.Status, &mandate.CreatedAt, &mandate.UpdatedAt,
		&mandate.CaseIDEncrypted, &mandate.ClientIDEncrypted,
		&mandate.MandateNumberEncrypted, &mandate.ScopeOfWorkEncrypted, &mandate.TermsConditionsEncrypted,
		&mandate.ClientSignatureEncrypted, &mandate.InvestigatorSignatureEncrypted, &mandate.SpecialInstructionsEncrypted,
		&mandate.IssueDate, &mandate.ValidFrom, &mandate.ValidUntil, &mandate.LinkedEstimateID, &mandate.Jurisdiction,
		&mandate.DEKEncrypted, &mandate.KeyVersion, &mandate.Metadata,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mandate not found")
		}
		return nil, fmt.Errorf("failed to get mandate: %w", err)
	}

	return &mandate, nil
}
