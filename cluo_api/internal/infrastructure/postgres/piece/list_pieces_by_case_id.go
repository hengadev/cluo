package pieceRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (r *Repository) ListPiecesByCaseID(ctx context.Context, caseID uuid.UUID, pagination piece.Pagination) ([]*piece.PieceEncx, int, error) {
	// Count total rows for this case
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s.pieces WHERE case_id = $1
	`, r.schema)

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, caseID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count pieces: %w", err)
	}

	// Paginated results
	query := fmt.Sprintf(`
		SELECT
			id, case_id, filename_encrypted, storage_key, mime_type, size_bytes,
			notes_encrypted, dek_encrypted, key_version, created_at, updated_at
		FROM %s.pieces
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, caseID, pagination.PageSize, pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pieces: %w", err)
	}
	defer rows.Close()

	var pieces []*piece.PieceEncx
	for rows.Next() {
		p := &piece.PieceEncx{}
		if err := rows.Scan(
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
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan piece row: %w", err)
		}
		pieces = append(pieces, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating piece rows: %w", err)
	}

	return pieces, total, nil
}
