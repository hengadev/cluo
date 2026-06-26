package aiTranscriptionJob

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

const Schema = "ai"

// Repository implements the TranscriptionJobRepository interface.
type Repository struct {
	pool     *pgxpool.Pool
	schema   string
	workerID string
}

// New creates a new transcription job repository.
func New(ctx context.Context, pool *pgxpool.Pool, workerID string) ports.TranscriptionJobRepository {
	return &Repository{
		pool:     pool,
		schema:   Schema,
		workerID: workerID,
	}
}

// Create creates a new transcription job.
func (r *Repository) Create(ctx context.Context, job *ai.TranscriptionJob) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.transcription_jobs (
			id, media_file_id, audio_path, status, progress,
			error_message, transcription_id, webhook_url,
			created_at, started_at, completed_at, claimed_at, claimed_by, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		job.ID,
		job.MediaFileID,
		job.AudioPath,
		job.Status,
		job.Progress,
		job.ErrorMessage,
		job.TranscriptionID,
		job.WebhookURL,
		job.CreatedAt,
		job.StartedAt,
		job.CompletedAt,
		job.ClaimedAt,
		job.ClaimedBy,
		job.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("create transcription job: %w", err)
	}

	return nil
}

// GetByID retrieves a job by ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*ai.TranscriptionJob, error) {
	query := fmt.Sprintf(`
		SELECT id, media_file_id, audio_path, status, progress,
		       error_message, transcription_id, webhook_url,
		       created_at, started_at, completed_at, claimed_at, claimed_by, created_by
		FROM %s.transcription_jobs
		WHERE id = $1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, id)

	var j ai.TranscriptionJob
	err := row.Scan(
		&j.ID, &j.MediaFileID, &j.AudioPath, &j.Status, &j.Progress,
		&j.ErrorMessage, &j.TranscriptionID, &j.WebhookURL,
		&j.CreatedAt, &j.StartedAt, &j.CompletedAt, &j.ClaimedAt, &j.ClaimedBy, &j.CreatedBy,
	)

	if err != nil {
		return nil, errs.ClassifyPgError("get transcription job by id", err)
	}

	return &j, nil
}

// GetByUserID retrieves jobs for a user with pagination.
func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ai.TranscriptionJob, error) {
	query := fmt.Sprintf(`
		SELECT id, media_file_id, audio_path, status, progress,
		       error_message, transcription_id, webhook_url,
		       created_at, started_at, completed_at, claimed_at, claimed_by, created_by
		FROM %s.transcription_jobs
		WHERE created_by = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query jobs by user id: %w", err)
	}
	defer rows.Close()

	var jobs []*ai.TranscriptionJob
	for rows.Next() {
		var j ai.TranscriptionJob
		err := rows.Scan(
			&j.ID, &j.MediaFileID, &j.AudioPath, &j.Status, &j.Progress,
			&j.ErrorMessage, &j.TranscriptionID, &j.WebhookURL,
			&j.CreatedAt, &j.StartedAt, &j.CompletedAt, &j.ClaimedAt, &j.ClaimedBy, &j.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transcription job: %w", err)
		}
		jobs = append(jobs, &j)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate transcription jobs: %w", rows.Err())
	}

	return jobs, nil
}

// GetByUserIDWithStatus retrieves jobs for a user with status filtering.
func (r *Repository) GetByUserIDWithStatus(ctx context.Context, userID uuid.UUID, status ai.JobStatus, limit, offset int) ([]*ai.TranscriptionJob, error) {
	query := fmt.Sprintf(`
		SELECT id, media_file_id, audio_path, status, progress,
		       error_message, transcription_id, webhook_url,
		       created_at, started_at, completed_at, claimed_at, claimed_by, created_by
		FROM %s.transcription_jobs
		WHERE created_by = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, userID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query jobs by user id with status: %w", err)
	}
	defer rows.Close()

	var jobs []*ai.TranscriptionJob
	for rows.Next() {
		var j ai.TranscriptionJob
		err := rows.Scan(
			&j.ID, &j.MediaFileID, &j.AudioPath, &j.Status, &j.Progress,
			&j.ErrorMessage, &j.TranscriptionID, &j.WebhookURL,
			&j.CreatedAt, &j.StartedAt, &j.CompletedAt, &j.ClaimedAt, &j.ClaimedBy, &j.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transcription job: %w", err)
		}
		jobs = append(jobs, &j)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate transcription jobs: %w", rows.Err())
	}

	return jobs, nil
}

// CountByUserID counts jobs for a user with optional status filter.
func (r *Repository) CountByUserID(ctx context.Context, userID uuid.UUID, status *ai.JobStatus) (int, error) {
	var query string
	var args []interface{}

	if status != nil {
		query = fmt.Sprintf(`SELECT COUNT(*) FROM %s.transcription_jobs WHERE created_by = $1 AND status = $2`, r.schema)
		args = []interface{}{userID, *status}
	} else {
		query = fmt.Sprintf(`SELECT COUNT(*) FROM %s.transcription_jobs WHERE created_by = $1`, r.schema)
		args = []interface{}{userID}
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count transcription jobs: %w", err)
	}

	return count, nil
}

// UpdateStatus updates the status of a job.
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status ai.JobStatus, errorMsg string) error {
	query := fmt.Sprintf(`
		UPDATE %s.transcription_jobs
		SET status = $2, error_message = $3
		WHERE id = $1
	`, r.schema)

	_, err := r.pool.Exec(ctx, query, id, status, errorMsg)
	if err != nil {
		return fmt.Errorf("update job status: %w", err)
	}

	return nil
}

// UpdateProgress updates the progress of a job.
func (r *Repository) UpdateProgress(ctx context.Context, id uuid.UUID, progress int) error {
	query := fmt.Sprintf(`
		UPDATE %s.transcription_jobs
		SET progress = $2
		WHERE id = $1
	`, r.schema)

	_, err := r.pool.Exec(ctx, query, id, progress)
	if err != nil {
		return fmt.Errorf("update job progress: %w", err)
	}

	return nil
}

// Complete marks a job as completed with a transcription ID.
func (r *Repository) Complete(ctx context.Context, id uuid.UUID, transcriptionID uuid.UUID) error {
	now := time.Now()
	query := fmt.Sprintf(`
		UPDATE %s.transcription_jobs
		SET status = $2,
		    transcription_id = $3,
		    progress = 100,
		    completed_at = $4
		WHERE id = $1
	`, r.schema)

	_, err := r.pool.Exec(ctx, query, id, ai.JobStatusCompleted, transcriptionID, now)
	if err != nil {
		return fmt.Errorf("complete job: %w", err)
	}

	return nil
}

// GetPendingJobs retrieves pending jobs up to the limit.
func (r *Repository) GetPendingJobs(ctx context.Context, limit int) ([]*ai.TranscriptionJob, error) {
	query := fmt.Sprintf(`
		SELECT id, media_file_id, audio_path, status, progress,
		       error_message, transcription_id, webhook_url,
		       created_at, started_at, completed_at, claimed_at, claimed_by, created_by
		FROM %s.transcription_jobs
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, ai.JobStatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("query pending jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*ai.TranscriptionJob
	for rows.Next() {
		var j ai.TranscriptionJob
		err := rows.Scan(
			&j.ID, &j.MediaFileID, &j.AudioPath, &j.Status, &j.Progress,
			&j.ErrorMessage, &j.TranscriptionID, &j.WebhookURL,
			&j.CreatedAt, &j.StartedAt, &j.CompletedAt, &j.ClaimedAt, &j.ClaimedBy, &j.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transcription job: %w", err)
		}
		jobs = append(jobs, &j)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate transcription jobs: %w", rows.Err())
	}

	return jobs, nil
}

// ClaimJob atomically claims a job for processing.
// Returns (claimed, error) - claimed is true if the job was successfully claimed.
func (r *Repository) ClaimJob(ctx context.Context, id uuid.UUID, workerID string) (bool, error) {
	now := time.Now()
	query := fmt.Sprintf(`
		UPDATE %s.transcription_jobs
		SET status = $2,
		    started_at = $3,
		    claimed_at = $3,
		    claimed_by = $4
		WHERE id = $1 AND status = $5
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id, ai.JobStatusProcessing, now, workerID, ai.JobStatusPending)
	if err != nil {
		return false, fmt.Errorf("claim job: %w", err)
	}

	rowsAffected := result.RowsAffected()
	return rowsAffected == 1, nil
}

// Delete deletes a job by ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.transcription_jobs WHERE id = $1`, r.schema)

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete job: %w", err)
	}

	return nil
}
