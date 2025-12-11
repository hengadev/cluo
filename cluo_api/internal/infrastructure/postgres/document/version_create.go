package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateVersion creates a new document version.
func (r *VersionRepository) CreateVersion(ctx context.Context, version *document.DocumentVersion) error {
	query := `
		INSERT INTO document_versions (
			id, document_id, doc_type, version, author_id, data, reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (document_id, doc_type, version) DO UPDATE SET
			author_id = EXCLUDED.author_id,
			data = EXCLUDED.data,
			reason = EXCLUDED.reason,
			created_at = EXCLUDED.created_at
	`

	_, err := r.pool.Exec(ctx, query,
		version.ID, version.DocumentID, version.DocType, version.Version,
		version.AuthorID, version.Data, version.Reason, version.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create document version: %w", err)
	}

	return nil
}
