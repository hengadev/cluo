package aiTranscriptAnalysis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

const Schema = "ai"

// Repository implements the TranscriptAnalysisRepository interface.
type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

// New creates a new transcript analysis repository.
func New(ctx context.Context, pool *pgxpool.Pool) ports.TranscriptAnalysisRepository {
	return &Repository{
		pool:   pool,
		schema: Schema,
	}
}

// Create creates or replaces the analysis for a transcription (encrypted).
// A transcription has at most one analysis, so re-analyzing replaces the previous result.
func (r *Repository) Create(ctx context.Context, analysis *ai.TranscriptAnalysisEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.transcript_analyses (
			id, transcription_id, key_findings_encrypted, summary_encrypted,
			sentiment, topics_encrypted, suggested_actions_encrypted,
			model_used, processing_time_ms, created_at, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		ON CONFLICT (transcription_id) DO UPDATE SET
			id = EXCLUDED.id,
			key_findings_encrypted = EXCLUDED.key_findings_encrypted,
			summary_encrypted = EXCLUDED.summary_encrypted,
			sentiment = EXCLUDED.sentiment,
			topics_encrypted = EXCLUDED.topics_encrypted,
			suggested_actions_encrypted = EXCLUDED.suggested_actions_encrypted,
			model_used = EXCLUDED.model_used,
			processing_time_ms = EXCLUDED.processing_time_ms,
			created_at = EXCLUDED.created_at,
			dek_encrypted = EXCLUDED.dek_encrypted,
			key_version = EXCLUDED.key_version
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		analysis.ID,
		analysis.TranscriptionID,
		analysis.KeyFindingsEncrypted,
		analysis.SummaryEncrypted,
		analysis.Sentiment,
		analysis.TopicsEncrypted,
		analysis.SuggestedActionsEncrypted,
		analysis.ModelUsed,
		analysis.ProcessingTimeMs,
		analysis.CreatedAt,
		analysis.DEKEncrypted,
		analysis.KeyVersion,
	)

	if err != nil {
		return fmt.Errorf("create transcript analysis: %w", err)
	}

	return nil
}

// GetByID retrieves an analysis by ID (encrypted).
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*ai.TranscriptAnalysisEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, transcription_id, key_findings_encrypted, summary_encrypted,
		       sentiment, topics_encrypted, suggested_actions_encrypted,
		       model_used, processing_time_ms, created_at, dek_encrypted, key_version
		FROM %s.transcript_analyses
		WHERE id = $1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, id)
	return scanAnalysis(row)
}

// GetByTranscriptionID retrieves an analysis by transcription ID (encrypted).
func (r *Repository) GetByTranscriptionID(ctx context.Context, transcriptionID uuid.UUID) (*ai.TranscriptAnalysisEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, transcription_id, key_findings_encrypted, summary_encrypted,
		       sentiment, topics_encrypted, suggested_actions_encrypted,
		       model_used, processing_time_ms, created_at, dek_encrypted, key_version
		FROM %s.transcript_analyses
		WHERE transcription_id = $1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, transcriptionID)
	return scanAnalysis(row)
}

// List retrieves analyses with pagination, optionally filtered by transcription ID (encrypted).
func (r *Repository) List(ctx context.Context, transcriptionID *uuid.UUID, limit, offset int) ([]*ai.TranscriptAnalysisEncx, int, error) {
	conditions := ""
	args := []any{}
	argIdx := 1

	if transcriptionID != nil {
		conditions = fmt.Sprintf("WHERE transcription_id = $%d", argIdx)
		args = append(args, *transcriptionID)
		argIdx++
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.transcript_analyses %s`, r.schema, conditions)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count transcript analyses: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT id, transcription_id, key_findings_encrypted, summary_encrypted,
		       sentiment, topics_encrypted, suggested_actions_encrypted,
		       model_used, processing_time_ms, created_at, dek_encrypted, key_version
		FROM %s.transcript_analyses
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, r.schema, conditions, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list transcript analyses: %w", err)
	}
	defer rows.Close()

	var results []*ai.TranscriptAnalysisEncx
	for rows.Next() {
		analysis, err := scanAnalysis(rows)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, analysis)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list transcript analyses: %w", err)
	}

	return results, total, nil
}

// Delete deletes an analysis by ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.transcript_analyses WHERE id = $1`, r.schema)

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete transcript analysis: %w", err)
	}

	return nil
}

// rowScanner abstracts pgx.Row and pgx.Rows for shared scanning logic.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanAnalysis(row rowScanner) (*ai.TranscriptAnalysisEncx, error) {
	var a ai.TranscriptAnalysisEncx
	err := row.Scan(
		&a.ID, &a.TranscriptionID, &a.KeyFindingsEncrypted, &a.SummaryEncrypted,
		&a.Sentiment, &a.TopicsEncrypted, &a.SuggestedActionsEncrypted,
		&a.ModelUsed, &a.ProcessingTimeMs, &a.CreatedAt, &a.DEKEncrypted, &a.KeyVersion,
	)
	if err != nil {
		return nil, fmt.Errorf("scan transcript analysis: %w", err)
	}
	return &a, nil
}
