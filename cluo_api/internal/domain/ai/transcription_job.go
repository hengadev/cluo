package ai

import (
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the status of a transcription job.
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

// Valid returns true if the job status is valid.
func (s JobStatus) Valid() bool {
	switch s {
	case JobStatusPending,
		JobStatusProcessing,
		JobStatusCompleted,
		JobStatusFailed,
		JobStatusCancelled:
		return true
	default:
		return false
	}
}

// IsTerminal returns true if the status is a terminal state.
func (s JobStatus) IsTerminal() bool {
	return s == JobStatusCompleted || s == JobStatusFailed || s == JobStatusCancelled
}

// IsCancellable returns true if the job can be cancelled.
func (s JobStatus) IsCancellable() bool {
	return s == JobStatusPending
}

// TranscriptionJob represents an asynchronous speech-to-text transcription job.
type TranscriptionJob struct {
	ID              uuid.UUID   `db:"id"`
	MediaFileID     uuid.UUID   `db:"media_file_id"`
	AudioPath       string      `db:"audio_path"`
	Status          JobStatus   `db:"status"`
	Progress        int         `db:"progress"`
	ErrorMessage    string      `db:"error_message"`
	TranscriptionID *uuid.UUID  `db:"transcription_id"`
	WebhookURL      *string     `db:"webhook_url"`
	CreatedAt       time.Time   `db:"created_at"`
	StartedAt       *time.Time  `db:"started_at"`
	CompletedAt     *time.Time  `db:"completed_at"`
	ClaimedAt       *time.Time  `db:"claimed_at"`
	ClaimedBy       *string     `db:"claimed_by"`
	CreatedBy       uuid.UUID   `db:"created_by"`
}

// NewTranscriptionJob creates a new TranscriptionJob entity.
func NewTranscriptionJob(
	mediaFileID uuid.UUID,
	audioPath string,
	webhookURL *string,
	createdBy uuid.UUID,
) *TranscriptionJob {
	return &TranscriptionJob{
		ID:          uuid.New(),
		MediaFileID: mediaFileID,
		AudioPath:   audioPath,
		Status:      JobStatusPending,
		Progress:    0,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		WebhookURL:  webhookURL,
	}
}

// Start marks the job as started.
func (j *TranscriptionJob) Start(workerID string) {
	j.Status = JobStatusProcessing
	now := time.Now()
	j.StartedAt = &now
	j.ClaimedAt = &now
	j.ClaimedBy = &workerID
}

// UpdateProgress updates the job progress percentage.
func (j *TranscriptionJob) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	j.Progress = progress
}

// Complete marks the job as completed with a transcription.
func (j *TranscriptionJob) Complete(transcriptionID uuid.UUID) {
	j.Status = JobStatusCompleted
	j.TranscriptionID = &transcriptionID
	j.Progress = 100
	now := time.Now()
	j.CompletedAt = &now
}

// Fail marks the job as failed with an error message.
func (j *TranscriptionJob) Fail(errorMessage string) {
	j.Status = JobStatusFailed
	j.ErrorMessage = errorMessage
	now := time.Now()
	j.CompletedAt = &now
}

// Cancel marks the job as cancelled.
func (j *TranscriptionJob) Cancel() {
	if j.Status.IsCancellable() {
		j.Status = JobStatusCancelled
		now := time.Now()
		j.CompletedAt = &now
	}
}
