package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// SpeechToTextService defines the interface for speech-to-text operations.
type SpeechToTextService interface {
	// SubmitTranscriptionJob submits a new transcription job for async processing.
	SubmitTranscriptionJob(ctx context.Context, req *SubmitTranscriptionRequest) (*ai.TranscriptionJob, error)

	// GetJobStatus retrieves the current status of a transcription job.
	GetJobStatus(ctx context.Context, jobID uuid.UUID) (*ai.TranscriptionJob, error)

	// CancelJob cancels a pending transcription job.
	CancelJob(ctx context.Context, jobID uuid.UUID) error

	// ListJobs retrieves jobs for a user with optional filtering.
	ListJobs(ctx context.Context, req *ListJobsRequest) ([]*ai.TranscriptionJob, int, error)

	// GetTranscription retrieves a completed transcription by ID.
	GetTranscription(ctx context.Context, id uuid.UUID) (*ai.Transcription, error)

	// GetTranscriptionByMediaFileID retrieves a transcription by media file ID.
	GetTranscriptionByMediaFileID(ctx context.Context, mediaFileID uuid.UUID) (*ai.Transcription, error)

	// DeleteTranscription deletes a transcription.
	DeleteTranscription(ctx context.Context, id uuid.UUID) error
}

// SubmitTranscriptionRequest defines the request for submitting a transcription job.
type SubmitTranscriptionRequest struct {
	MediaFileID uuid.UUID
	AudioData   []byte
	AudioFilename string
	WebhookURL  *string
	UserID      uuid.UUID
}

// ListJobsRequest defines the request for listing transcription jobs.
type ListJobsRequest struct {
	UserID  uuid.UUID
	Status  *ai.JobStatus
	Limit   int
	Offset  int
}
