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
			title_encrypted = $4,
			description_encrypted = $5,
			status_encrypted = $6,
			updated_at_encrypted = $7,
			dek_encrypted = $8,
			key_version = $9,
			metadata = $10
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.TitleEncrypted,
		caseEncx.DescriptionEncrypted,
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

