package mediaRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (r *Repository) CreateMedia(ctx context.Context, mediaEncx *domain.MediaFileEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.media_files (
			id, caseid, filesize, ispublished, createdat,
			url_encrypted, type_encrypted, mimetype_encrypted,
			filename_encrypted, caption_encrypted,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		mediaEncx.ID,
		mediaEncx.CaseID,
		mediaEncx.FileSize,
		mediaEncx.IsPublished,
		mediaEncx.CreatedAt,
		mediaEncx.URLEncrypted,
		mediaEncx.TypeEncrypted,
		mediaEncx.MimeTypeEncrypted,
		mediaEncx.FileNameEncrypted,
		mediaEncx.CaptionEncrypted,
		mediaEncx.DEKEncrypted,
		mediaEncx.KeyVersion,
		mediaEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create media", err)
	}

	return nil
}
