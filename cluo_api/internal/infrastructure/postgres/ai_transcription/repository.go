package aiTranscription

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

const Schema = "ai"

// Repository implements the TranscriptionRepository interface.
type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

// New creates a new transcription repository.
func New(ctx context.Context, pool *pgxpool.Pool) ports.TranscriptionRepository {
	return &Repository{
		pool:   pool,
		schema: Schema,
	}
}

// Create creates a new transcription (encrypted).
func (r *Repository) Create(ctx context.Context, transcription *ai.TranscriptionEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.transcriptions (
			id, job_id, media_file_id, audio_url, transcript_encrypted,
			confidence_score, language, duration, model_name, processing_time_ms,
			created_at, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		transcription.ID,
		transcription.JobID,
		transcription.MediaFileID,
		transcription.AudioURL,
		transcription.TranscriptEncrypted,
		transcription.ConfidenceScore,
		transcription.Language,
		transcription.Duration,
		transcription.ModelName,
		transcription.ProcessingTimeMs,
		transcription.CreatedAt,
		transcription.DEKEncrypted,
		transcription.KeyVersion,
	)

	if err != nil {
		return fmt.Errorf("create transcription: %w", err)
	}

	return nil
}

// GetByID retrieves a transcription by ID (encrypted).
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*ai.TranscriptionEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, job_id, media_file_id, audio_url, transcript_encrypted,
		       confidence_score, language, duration, model_name, processing_time_ms,
		       created_at, audio_deleted_at, dek_encrypted, key_version
		FROM %s.transcriptions
		WHERE id = $1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, id)

	var t ai.TranscriptionEncx
	err := row.Scan(
		&t.ID, &t.JobID, &t.MediaFileID, &t.AudioURL, &t.TranscriptEncrypted,
		&t.ConfidenceScore, &t.Language, &t.Duration, &t.ModelName, &t.ProcessingTimeMs,
		&t.CreatedAt, &t.AudioDeletedAt, &t.DEKEncrypted, &t.KeyVersion,
	)

	if err != nil {
		return nil, fmt.Errorf("get transcription by id: %w", err)
	}

	return &t, nil
}

// GetByMediaFileID retrieves a transcription by media file ID (encrypted).
func (r *Repository) GetByMediaFileID(ctx context.Context, mediaFileID uuid.UUID) (*ai.TranscriptionEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, job_id, media_file_id, audio_url, transcript_encrypted,
		       confidence_score, language, duration, model_name, processing_time_ms,
		       created_at, audio_deleted_at, dek_encrypted, key_version
		FROM %s.transcriptions
		WHERE media_file_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, mediaFileID)

	var t ai.TranscriptionEncx
	err := row.Scan(
		&t.ID, &t.JobID, &t.MediaFileID, &t.AudioURL, &t.TranscriptEncrypted,
		&t.ConfidenceScore, &t.Language, &t.Duration, &t.ModelName, &t.ProcessingTimeMs,
		&t.CreatedAt, &t.AudioDeletedAt, &t.DEKEncrypted, &t.KeyVersion,
	)

	if err != nil {
		return nil, fmt.Errorf("get transcription by media file id: %w", err)
	}

	return &t, nil
}

// GetByJobID retrieves a transcription by job ID (encrypted).
func (r *Repository) GetByJobID(ctx context.Context, jobID uuid.UUID) (*ai.TranscriptionEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, job_id, media_file_id, audio_url, transcript_encrypted,
		       confidence_score, language, duration, model_name, processing_time_ms,
		       created_at, audio_deleted_at, dek_encrypted, key_version
		FROM %s.transcriptions
		WHERE job_id = $1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, jobID)

	var t ai.TranscriptionEncx
	err := row.Scan(
		&t.ID, &t.JobID, &t.MediaFileID, &t.AudioURL, &t.TranscriptEncrypted,
		&t.ConfidenceScore, &t.Language, &t.Duration, &t.ModelName, &t.ProcessingTimeMs,
		&t.CreatedAt, &t.AudioDeletedAt, &t.DEKEncrypted, &t.KeyVersion,
	)

	if err != nil {
		return nil, fmt.Errorf("get transcription by job id: %w", err)
	}

	return &t, nil
}

// Delete deletes a transcription by ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.transcriptions WHERE id = $1`, r.schema)

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete transcription: %w", err)
	}

	return nil
}
