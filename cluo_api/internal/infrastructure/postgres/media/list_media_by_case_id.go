package mediaRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (r *Repository) ListMediaByCaseID(
	ctx context.Context,
	caseID uuid.UUID,
	mediaType *domain.MediaType,
	page, pageSize int,
) ([]*domain.MediaFileEncx, int, error) {
	// Note: Type is encrypted, so we can't filter at DB level
	// Service layer will filter after decryption if mediaType is provided

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.media_files
		WHERE caseid = $1
	`, r.schema)

	var total int
	err := r.pool.QueryRow(ctx, countQuery, caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count media: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT
			id, caseid, filesize, ispublished, purpose, createdat,
			url_encrypted, type_encrypted, mimetype_encrypted,
			filename_encrypted, caption_encrypted,
			dek_encrypted, key_version, metadata
		FROM %s.media_files
		WHERE caseid = $1
		ORDER BY createdat DESC
		LIMIT $2 OFFSET $3
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, caseID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query media: %w", err)
	}
	defer rows.Close()

	var mediaList []*domain.MediaFileEncx
	for rows.Next() {
		mediaEncx := &domain.MediaFileEncx{}
		err := rows.Scan(
			&mediaEncx.ID,
			&mediaEncx.CaseID,
			&mediaEncx.FileSize,
			&mediaEncx.IsPublished,
			&mediaEncx.Purpose,
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
			return nil, 0, fmt.Errorf("failed to scan media row: %w", err)
		}

		mediaList = append(mediaList, mediaEncx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating media rows: %w", err)
	}

	return mediaList, total, nil
}
