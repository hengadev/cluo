package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// DeleteVersions deletes all versions for a document.
func (r *VersionRepository) DeleteVersions(ctx context.Context, documentID string, docType document.DocumentType) error {
	query := `DELETE FROM document_versions WHERE document_id = $1 AND doc_type = $2`
	_, err := r.pool.Exec(ctx, query, documentID, docType)
	if err != nil {
		return fmt.Errorf("failed to delete document versions: %w", err)
	}
	return nil
}