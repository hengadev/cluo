package mediaRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (r *Repository) UpdateMedia(ctx context.Context, mediaEncx *domain.MediaFileEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.media_files
		SET
			caption_encrypted = $2,
			ispublished = $3,
			dek_encrypted = $4,
			key_version = $5,
			metadata = $6
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		mediaEncx.ID,
		mediaEncx.CaptionEncrypted,
		mediaEncx.IsPublished,
		mediaEncx.DEKEncrypted,
		mediaEncx.KeyVersion,
		mediaEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("update media", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("media with ID %s not found", mediaEncx.ID), "media")
	}

	return nil
}
