package caseRepository

import (
	"context"
	"fmt"
	"time"

	caseDomain "github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/case"
)

func (r *Repository) CreateCase(ctx context.Context, caseEncx *caseDomain.CaseEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.cases (
			id, clientid, assignedcontactid, createdat,
			title_encrypted, description_encrypted, status_encrypted,
			updatedat_encrypted, dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		caseEncx.ID,
		caseEncx.ClientID,
		caseEncx.AssignedContactID,
		caseEncx.CreatedAt,
		caseEncx.TitleEncrypted,
		caseEncx.DescriptionEncrypted,
		caseEncx.StatusEncrypted,
		caseEncx.UpdatedAtEncrypted,
		caseEncx.DEKEncrypted,
		caseEncx.KeyVersion,
		caseEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create case", err)
	}

	return nil
}