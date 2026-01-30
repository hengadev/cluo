package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// TranscriptionJobRepository defines the interface for transcription job persistence.
type TranscriptionJobRepository interface {
	// Create creates a new transcription job.
	Create(ctx context.Context, job *ai.TranscriptionJob) error

	// GetByID retrieves a job by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*ai.TranscriptionJob, error)

	// GetByUserID retrieves jobs for a user with pagination.
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ai.TranscriptionJob, error)

	// GetByUserIDWithStatus retrieves jobs for a user with status filtering.
	GetByUserIDWithStatus(ctx context.Context, userID uuid.UUID, status ai.JobStatus, limit, offset int) ([]*ai.TranscriptionJob, error)

	// CountByUserID counts jobs for a user with optional status filter.
	CountByUserID(ctx context.Context, userID uuid.UUID, status *ai.JobStatus) (int, error)

	// UpdateStatus updates the status of a job.
	UpdateStatus(ctx context.Context, id uuid.UUID, status ai.JobStatus, errorMsg string) error

	// UpdateProgress updates the progress of a job.
	UpdateProgress(ctx context.Context, id uuid.UUID, progress int) error

	// Complete marks a job as completed with a transcription ID.
	Complete(ctx context.Context, id uuid.UUID, transcriptionID uuid.UUID) error

	// GetPendingJobs retrieves pending jobs up to the limit.
	GetPendingJobs(ctx context.Context, limit int) ([]*ai.TranscriptionJob, error)

	// ClaimJob atomically claims a job for processing.
	// Returns (claimed, error) - claimed is true if the job was successfully claimed.
	ClaimJob(ctx context.Context, id uuid.UUID, workerID string) (bool, error)

	// Delete deletes a job by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
