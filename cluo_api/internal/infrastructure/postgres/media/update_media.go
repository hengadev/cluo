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
			url_encrypted = $2,
			type_encrypted = $3,
			mimetype_encrypted = $4,
			filename_encrypted = $5,
			caption_encrypted = $6,
			ispublished = $7,
			purpose = $8,
			dek_encrypted = $9,
			key_version = $10,
			metadata = $11
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		mediaEncx.ID,
		mediaEncx.URLEncrypted,
		mediaEncx.TypeEncrypted,
		mediaEncx.MimeTypeEncrypted,
		mediaEncx.FileNameEncrypted,
		mediaEncx.CaptionEncrypted,
		mediaEncx.IsPublished,
		mediaEncx.Purpose,
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
