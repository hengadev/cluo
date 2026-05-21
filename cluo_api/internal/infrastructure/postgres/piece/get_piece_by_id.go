package pieceRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (r *Repository) GetPieceByID(ctx context.Context, id uuid.UUID) (*piece.PieceEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, case_id, filename_encrypted, storage_key, mime_type, size_bytes,
			notes_encrypted, dek_encrypted, key_version, created_at, updated_at
		FROM %s.pieces
		WHERE id = $1
	`, r.schema)

	p := &piece.PieceEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.CaseID,
		&p.FilenameEncrypted,
		&p.StorageKey,
		&p.MimeType,
		&p.SizeBytes,
		&p.NotesEncrypted,
		&p.DEKEncrypted,
		&p.KeyVersion,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get piece by id", err)
	}

	return p, nil
}
