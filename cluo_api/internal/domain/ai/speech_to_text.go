package ai

import (
	"time"

	"github.com/google/uuid"
)

// Transcription represents a speech-to-text transcription result.
// Note: DEKEncrypted and KeyVersion are auto-generated in TranscriptionEncx by encx.
type Transcription struct {
	ID               uuid.UUID  `db:"id"`
	JobID            uuid.UUID  `db:"job_id"`
	MediaFileID      uuid.UUID  `db:"media_file_id"`
	AudioURL         string     `db:"audio_url"`
	Transcript       string     `encx:"encrypt" db:"transcript_encrypted"`
	ConfidenceScore  float32    `db:"confidence_score"`
	Language         string     `db:"language"`
	Duration         int64      `db:"duration"` // in milliseconds
	ModelName        string     `db:"model_name"`
	ProcessingTimeMs int64      `db:"processing_time_ms"`
	CreatedAt        time.Time  `db:"created_at"`
	AudioDeletedAt   *time.Time `db:"audio_deleted_at"`
}

// NewTranscription creates a new Transcription entity.
func NewTranscription(
	jobID, mediaFileID uuid.UUID,
	audioURL, transcript string,
	confidenceScore float32,
	language string,
	duration int64,
	modelName string,
	processingTimeMs int64,
) *Transcription {
	return &Transcription{
		ID:               uuid.New(),
		JobID:            jobID,
		MediaFileID:      mediaFileID,
		AudioURL:         audioURL,
		Transcript:       transcript,
		ConfidenceScore:  confidenceScore,
		Language:         language,
		Duration:         duration,
		ModelName:        modelName,
		ProcessingTimeMs: processingTimeMs,
		CreatedAt:        time.Now(),
	}
}

// MarkAudioDeleted marks the audio file as deleted.
func (t *Transcription) MarkAudioDeleted() {
	now := time.Now()
	t.AudioDeletedAt = &now
}
