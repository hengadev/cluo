package mediaHelpers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearMediaTable removes all media from database
func ClearMediaTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(ctx, "TRUNCATE TABLE media.media_files CASCADE")
	require.NoError(t, err, "Failed to clear media_files table")
}

// InsertMediaEncx inserts a MediaFileEncx into database
func InsertMediaEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool, mediaEncx *domain.MediaFileEncx) {
	t.Helper()

	query := `
		INSERT INTO media.media_files (
			id, caseid, filesize, ispublished, createdat,
			url_encrypted, type_encrypted, mimetype_encrypted,
			filename_encrypted, caption_encrypted,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := pool.Exec(ctx, query,
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
	require.NoError(t, err, "Failed to insert media")
}

// GetMediaEncxByID retrieves a MediaFileEncx from database
func GetMediaEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) *domain.MediaFileEncx {
	t.Helper()

	query := `
		SELECT
			id, caseid, filesize, ispublished, createdat,
			url_encrypted, type_encrypted, mimetype_encrypted,
			filename_encrypted, caption_encrypted,
			dek_encrypted, key_version, metadata
		FROM media.media_files
		WHERE id = $1
	`

	mediaEncx := &domain.MediaFileEncx{}
	err := pool.QueryRow(ctx, query, id).Scan(
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
	require.NoError(t, err, "Failed to get media")

	return mediaEncx
}

// CountMediaByCaseID counts media files for a case
func CountMediaByCaseID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID uuid.UUID) int {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM media.media_files WHERE caseid = $1",
		caseID).Scan(&count)
	require.NoError(t, err, "Failed to count media")

	return count
}
