package documentRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VersionRepository struct {
	pool *pgxpool.Pool
}

func NewVersionRepository(pool *pgxpool.Pool) ports.DocumentVersionRepository {
	return &VersionRepository{pool: pool}
}

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

func (r *VersionRepository) GetVersion(ctx context.Context, documentID string, docType document.DocumentType, version int) (*document.DocumentVersion, error) {
	query := `
		SELECT id, document_id, doc_type, version, author_id, data, reason, created_at
		FROM document_versions
		WHERE document_id = $1 AND doc_type = $2 AND version = $3
	`

	row := r.pool.QueryRow(ctx, query, documentID, docType, version)
	var documentVersion document.DocumentVersion

	err := row.Scan(
		&documentVersion.ID, &documentVersion.DocumentID, &documentVersion.DocType,
		&documentVersion.Version, &documentVersion.AuthorID, &documentVersion.Data,
		&documentVersion.Reason, &documentVersion.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document version not found")
		}
		return nil, fmt.Errorf("failed to get document version: %w", err)
	}

	return &documentVersion, nil
}

func (r *VersionRepository) DeleteVersions(ctx context.Context, documentID string, docType document.DocumentType) error {
	query := `DELETE FROM document_versions WHERE document_id = $1 AND doc_type = $2`
	_, err := r.pool.Exec(ctx, query, documentID, docType)
	if err != nil {
		return fmt.Errorf("failed to delete document versions: %w", err)
	}
	return nil
}
