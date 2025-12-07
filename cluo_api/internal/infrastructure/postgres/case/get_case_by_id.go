package caseRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (r *Repository) GetCaseByID(ctx context.Context, caseID uuid.UUID) (*caseDomain.CaseEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, client_id, assigned_contact_id, created_at,
			title_encrypted, description_encrypted, status_encrypted,
			updated_at_encrypted, dek_encrypted, key_version, metadata
		FROM %s.cases WHERE id = $1
	`, r.schema)

	caseEncx := &caseDomain.CaseEncx{}

	err := r.pool.QueryRow(ctx, query, caseID).Scan(
		&caseEncx.ID, &caseEncx.ClientID, &caseEncx.AssignedContactID, &caseEncx.CreatedAt,
		&caseEncx.TitleEncrypted, &caseEncx.DescriptionEncrypted, &caseEncx.StatusEncrypted,
		&caseEncx.UpdatedAtEncrypted, &caseEncx.DEKEncrypted, &caseEncx.KeyVersion, &caseEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get case by id", err)
	}
	return caseEncx, nil
}
