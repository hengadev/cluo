package mediaRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.media_files
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete media", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("media with ID %s not found", id), "media")
	}

	return nil
}
