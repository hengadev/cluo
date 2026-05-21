package rapportRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (r *Repository) UpdateRapport(ctx context.Context, rEncx *rapport.RapportEncx) error {
	if rEncx == nil {
		return fmt.Errorf("rapport cannot be nil")
	}

	query := fmt.Sprintf(`
		UPDATE %s.rapports SET
			content_encrypted = $2,
			dek_encrypted = $3,
			key_version = $4,
			updated_at = $5
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		rEncx.ID,
		rEncx.ContentEncrypted,
		rEncx.DEKEncrypted,
		rEncx.KeyVersion,
		rEncx.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("update rapport", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "rapport for update")
	}

	return nil
}
