package rapportRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (r *Repository) CreateRapport(ctx context.Context, rEncx *rapport.RapportEncx) error {
	if rEncx == nil {
		return fmt.Errorf("rapport cannot be nil")
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.rapports (
			id, case_id, content_encrypted, dek_encrypted, key_version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		rEncx.ID,
		rEncx.CaseID,
		rEncx.ContentEncrypted,
		rEncx.DEKEncrypted,
		rEncx.KeyVersion,
		rEncx.CreatedAt,
		rEncx.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create rapport", err)
	}

	return nil
}
