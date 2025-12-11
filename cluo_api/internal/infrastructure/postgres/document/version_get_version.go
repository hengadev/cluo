package documentRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetVersion retrieves a specific version of a document.
func (r *VersionRepository) GetVersion(ctx context.Context, documentID string, docType document.DocumentType, version int) (*document.DocumentVersion, error) {
	query := `
		SELECT id, document_id, doc_type, version, author_id, data, reason, created_at
		FROM document_versions
		WHERE document_id = $1 AND doc_type = $2 AND version = $3
	`

	row := r.pool.QueryRow(ctx, query, documentID, docType, version)
	var docVersion document.DocumentVersion

	err := row.Scan(
		&docVersion.ID, &docVersion.DocumentID, &docVersion.DocType, &docVersion.Version,
		&docVersion.AuthorID, &docVersion.Data, &docVersion.Reason, &docVersion.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document version not found")
		}
		return nil, fmt.Errorf("failed to get document version: %w", err)
	}

	return &docVersion, nil
}