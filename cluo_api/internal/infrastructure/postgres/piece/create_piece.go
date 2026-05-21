package pieceRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (r *Repository) CreatePiece(ctx context.Context, p *piece.PieceEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.pieces (
			id, case_id, filename_encrypted, storage_key, mime_type, size_bytes,
			notes_encrypted, dek_encrypted, key_version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		p.ID,
		p.CaseID,
		p.FilenameEncrypted,
		p.StorageKey,
		p.MimeType,
		p.SizeBytes,
		p.NotesEncrypted,
		p.DEKEncrypted,
		p.KeyVersion,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create piece", err)
	}

	return nil
}
