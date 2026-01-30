package aiSpeechToText

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/infrastructure/whisper"
	"github.com/hengadev/cluo_api/internal/infrastructure/worker"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

// Service implements the SpeechToTextService interface.
type Service struct {
	jobRepo           ports.TranscriptionJobRepository
	transcriptionRepo ports.TranscriptionRepository
	worker            *worker.TranscriptionWorker
	audioProcessor    *whisper.AudioProcessor
	crypto            encx.CryptoService
	tempAudioDir      string
	logger            *slog.Logger
}

// New creates a new speech-to-text service.
func New(
	jobRepo ports.TranscriptionJobRepository,
	transcriptionRepo ports.TranscriptionRepository,
	worker *worker.TranscriptionWorker,
	audioProcessor *whisper.AudioProcessor,
	crypto encx.CryptoService,
	tempAudioDir string,
	logger *slog.Logger,
) *Service {
	return &Service{
		jobRepo:           jobRepo,
		transcriptionRepo: transcriptionRepo,
		worker:            worker,
		audioProcessor:    audioProcessor,
		crypto:            crypto,
		tempAudioDir:      tempAudioDir,
		logger:            logger.With("component", "ai_speech_to_text"),
	}
}

// SubmitTranscriptionJob submits a new transcription job for async processing.
func (s *Service) SubmitTranscriptionJob(ctx context.Context, req *ports.SubmitTranscriptionRequest) (*ai.TranscriptionJob, error) {
	// Validate audio data
	if len(req.AudioData) == 0 {
		return nil, fmt.Errorf("audio data is empty")
	}

	// Save audio to temporary storage
	audioPath, err := s.saveTempAudio(req.AudioData, req.AudioFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to save audio: %w", err)
	}

	// Create job record
	job := ai.NewTranscriptionJob(
		req.MediaFileID,
		audioPath,
		req.WebhookURL,
		req.UserID,
	)

	if err := s.jobRepo.Create(ctx, job); err != nil {
		// Cleanup on failure
		os.Remove(audioPath)
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	s.logger.Info("transcription job submitted",
		"jobID", job.ID,
		"mediaFileID", job.MediaFileID,
		"userID", job.CreatedBy)

	// Enqueue for processing
	s.worker.Enqueue(job.ID)

	return job, nil
}

// GetJobStatus retrieves the current status of a transcription job.
func (s *Service) GetJobStatus(ctx context.Context, jobID uuid.UUID) (*ai.TranscriptionJob, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

// CancelJob cancels a pending transcription job.
func (s *Service) CancelJob(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	if !job.Status.IsCancellable() {
		return fmt.Errorf("%w: current status is %s", ai.ErrAIJobNotCancellable, job.Status)
	}

	// Update job status
	if err := s.jobRepo.UpdateStatus(ctx, jobID, ai.JobStatusCancelled, "cancelled by user"); err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}

	// Cleanup audio file
	if err := os.Remove(job.AudioPath); err != nil && !os.IsNotExist(err) {
		s.logger.Warn("failed to delete audio file after cancellation",
			"jobID", jobID,
			"error", err)
	}

	s.logger.Info("transcription job cancelled", "jobID", jobID)

	return nil
}

// ListJobs retrieves jobs for a user with optional filtering.
func (s *Service) ListJobs(ctx context.Context, req *ports.ListJobsRequest) ([]*ai.TranscriptionJob, int, error) {
	// Apply defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	var jobs []*ai.TranscriptionJob
	var err error

	if req.Status != nil {
		jobs, err = s.jobRepo.GetByUserIDWithStatus(ctx, req.UserID, *req.Status, req.Limit, req.Offset)
	} else {
		jobs, err = s.jobRepo.GetByUserID(ctx, req.UserID, req.Limit, req.Offset)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list jobs: %w", err)
	}

	// Get total count
	total, err := s.jobRepo.CountByUserID(ctx, req.UserID, req.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	return jobs, total, nil
}

// GetTranscription retrieves a completed transcription by ID.
func (s *Service) GetTranscription(ctx context.Context, id uuid.UUID) (*ai.Transcription, error) {
	// Get encrypted transcription from repository
	transcriptionEncx, err := s.transcriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcription: %w", err)
	}

	// Decrypt using encx
	transcription, err := ai.DecryptTranscriptionEncx(ctx, s.crypto, transcriptionEncx)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transcription: %w", err)
	}

	return transcription, nil
}

// GetTranscriptionByMediaFileID retrieves a transcription by media file ID.
func (s *Service) GetTranscriptionByMediaFileID(ctx context.Context, mediaFileID uuid.UUID) (*ai.Transcription, error) {
	// Get encrypted transcription from repository
	transcriptionEncx, err := s.transcriptionRepo.GetByMediaFileID(ctx, mediaFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcription: %w", err)
	}

	// Decrypt using encx
	transcription, err := ai.DecryptTranscriptionEncx(ctx, s.crypto, transcriptionEncx)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transcription: %w", err)
	}

	return transcription, nil
}

// DeleteTranscription deletes a transcription.
func (s *Service) DeleteTranscription(ctx context.Context, id uuid.UUID) error {
	if err := s.transcriptionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete transcription: %w", err)
	}

	s.logger.Info("transcription deleted", "id", id)

	return nil
}

// saveTempAudio saves audio data to a temporary file.
func (s *Service) saveTempAudio(data []byte, filename string) (string, error) {
	// Ensure temp directory exists
	if err := os.MkdirAll(s.tempAudioDir, 0755); err != nil {
		return "", fmt.Errorf("create temp directory: %w", err)
	}

	// Generate unique filename with timestamp and UUID
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".wav" // Default to wav
	}

	tempFilename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
	tempPath := filepath.Join(s.tempAudioDir, tempFilename)

	// Write audio data
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", fmt.Errorf("write audio file: %w", err)
	}

	return tempPath, nil
}

// HealthCheck checks if the speech-to-text service is healthy.
func (s *Service) HealthCheck(ctx context.Context) error {
	// Check if worker is running
	if s.worker == nil {
		return fmt.Errorf("worker not initialized")
	}

	return nil
}
