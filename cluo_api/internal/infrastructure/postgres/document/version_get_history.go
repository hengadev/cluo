package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetDocumentHistory retrieves the version history for a document.
func (r *VersionRepository) GetDocumentHistory(ctx context.Context, documentID string, docType document.DocumentType, pagination document.Pagination) ([]*document.DocumentVersion, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM document_versions WHERE document_id = $1 AND doc_type = $2",
		documentID, docType,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count document versions: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize

	query := `
		SELECT id, document_id, doc_type, version, author_id, data, reason, created_at
		FROM document_versions
		WHERE document_id = $1 AND doc_type = $2
		ORDER BY version DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, documentID, docType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query document versions: %w", err)
	}
	defer rows.Close()

	var versions []*document.DocumentVersion
	for rows.Next() {
		var version document.DocumentVersion
		err := rows.Scan(
			&version.ID, &version.DocumentID, &version.DocType, &version.Version,
			&version.AuthorID, &version.Data, &version.Reason, &version.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan document version: %w", err)
		}
		versions = append(versions, &version)
	}

	return versions, total, nil
}