package mediaRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (r *Repository) GetMediaByID(ctx context.Context, id uuid.UUID) (*domain.MediaFileEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, caseid, filesize, ispublished, createdat,
			url_encrypted, type_encrypted, mimetype_encrypted,
			filename_encrypted, caption_encrypted,
			dek_encrypted, key_version, metadata
		FROM %s.media_files
		WHERE id = $1
	`, r.schema)

	mediaEncx := &domain.MediaFileEncx{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&mediaEncx.ID,
		&mediaEncx.CaseID,
		&mediaEncx.FileSize,
		&mediaEncx.IsPublished,
		&mediaEncx.CreatedAt,
		&mediaEncx.URLEncrypted,
		&mediaEncx.TypeEncrypted,
		&mediaEncx.MimeTypeEncrypted,
		&mediaEncx.FileNameEncrypted,
		&mediaEncx.CaptionEncrypted,
		&mediaEncx.DEKEncrypted,
		&mediaEncx.KeyVersion,
		&mediaEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get media by id", err)
	}

	return mediaEncx, nil
}
