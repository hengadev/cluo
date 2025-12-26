package caseRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (r *Repository) UpdateCase(ctx context.Context, caseEncx *caseDomain.CaseEncx) error {
	if caseEncx == nil {
		return fmt.Errorf("case cannot be nil")
	}

	query := fmt.Sprintf(`
		UPDATE %s.cases SET
			client_id = $2,
			assigned_contact_id = $3,
			case_type = $4,
			title_encrypted = $5,
			description_encrypted = $6,
			external_reference_encrypted = $7,
			status_encrypted = $8,
			updated_at_encrypted = $9,
			dek_encrypted = $10,
			key_version = $11,
			metadata = $12
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.CaseType,
		caseEncx.TitleEncrypted,
		caseEncx.DescriptionEncrypted,
		caseEncx.ExternalReferenceEncrypted,
		caseEncx.StatusEncrypted,
		caseEncx.UpdatedAtEncrypted,
		caseEncx.DEKEncrypted,
		caseEncx.KeyVersion,
		caseEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("update case", err)
	}

	// Check if any row was actually updated
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "case for update")
	}

	return nil
}

