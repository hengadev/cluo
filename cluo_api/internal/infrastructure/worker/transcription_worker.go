package worker

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/infrastructure/webhook"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

const (
	// DefaultPollInterval is the default interval for polling pending jobs.
	DefaultPollInterval = 30 * time.Second
	// WorkerIDPrefix is the prefix for worker IDs.
	WorkerIDPrefix = "worker-"
)

// TranscriptionWorker processes transcription jobs in the background.
type TranscriptionWorker struct {
	whisperClient     ports.WhisperClient
	jobRepo           ports.TranscriptionJobRepository
	transcriptionRepo ports.TranscriptionRepository
	mediaRepo         ports.MediaRepository
	webhookClient     *webhook.Client
	crypto            encx.CryptoService
	queue             *JobQueue
	concurrency       int
	deleteAudioAfter  bool
	workerID          string
	stopCh            chan struct{}
	wg                sync.WaitGroup
	logger            *slog.Logger
}

// NewTranscriptionWorker creates a new transcription worker.
func NewTranscriptionWorker(
	whisperClient ports.WhisperClient,
	jobRepo ports.TranscriptionJobRepository,
	transcriptionRepo ports.TranscriptionRepository,
	mediaRepo ports.MediaRepository,
	webhookClient *webhook.Client,
	crypto encx.CryptoService,
	concurrency int,
	queueSize int,
	deleteAudioAfter bool,
	logger *slog.Logger,
) *TranscriptionWorker {
	return &TranscriptionWorker{
		whisperClient:     whisperClient,
		jobRepo:           jobRepo,
		transcriptionRepo: transcriptionRepo,
		mediaRepo:         mediaRepo,
		webhookClient:     webhookClient,
		crypto:            crypto,
		queue:             NewJobQueue(queueSize, logger),
		concurrency:       concurrency,
		deleteAudioAfter:  deleteAudioAfter,
		workerID:          WorkerIDPrefix + uuid.New().String()[:8],
		stopCh:            make(chan struct{}),
		logger:            logger.With("component", "transcription_worker"),
	}
}

// Start starts the worker goroutines.
func (w *TranscriptionWorker) Start(ctx context.Context) {
	w.logger.Info("starting transcription worker",
		"workerID", w.workerID,
		"concurrency", w.concurrency)

	// Start worker goroutines
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)
		go w.processLoop(ctx, i)
	}

	// Start pending job polling goroutine
	w.wg.Add(1)
	go w.pollPendingJobs(ctx)

	w.logger.Info("transcription worker started", "workerID", w.workerID)
}

// Stop stops the worker gracefully.
func (w *TranscriptionWorker) Stop() {
	w.logger.Info("stopping transcription worker", "workerID", w.workerID)

	close(w.stopCh)
	w.wg.Wait()

	w.logger.Info("transcription worker stopped", "workerID", w.workerID)
}

// Enqueue adds a job to the processing queue.
func (w *TranscriptionWorker) Enqueue(jobID uuid.UUID) {
	if !w.queue.Enqueue(jobID) {
		w.logger.Warn("job queue full", "jobID", jobID)
	}
}

// processLoop is the main processing loop for a worker goroutine.
func (w *TranscriptionWorker) processLoop(ctx context.Context, workerNum int) {
	defer w.wg.Done()

	workerLogger := w.logger.With("worker", workerNum)
	workerLogger.Info("worker processing loop started")

	for {
		select {
		case <-w.stopCh:
			workerLogger.Debug("worker received stop signal")
			return
		case <-ctx.Done():
			workerLogger.Debug("worker context cancelled")
			return
		default:
		}

		// Dequeue a job (with timeout for periodic stop checks)
		dequeueCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		jobID, ok := w.queue.Dequeue(dequeueCtx)
		cancel()

		if !ok {
			// Context cancelled or stopped, check exit conditions
			select {
			case <-w.stopCh:
				return
			case <-ctx.Done():
				return
			default:
				// Timeout, continue loop
				continue
			}
		}

		// Process the job
		w.processJob(ctx, jobID, workerNum)
	}
}

// processJob processes a single transcription job.
func (w *TranscriptionWorker) processJob(ctx context.Context, jobID uuid.UUID, workerNum int) {
	workerLogger := w.logger.With("jobID", jobID, "worker", workerNum)
	workerLogger.Info("processing job")

	// 1. Claim job atomically (prevents duplicate processing)
	claimed, err := w.jobRepo.ClaimJob(ctx, jobID, w.workerID)
	if err != nil {
		workerLogger.Error("failed to claim job", "error", err)
		return
	}

	if !claimed {
		workerLogger.Debug("job already claimed by another worker")
		return
	}

	// 2. Update status to processing
	if err := w.jobRepo.UpdateStatus(ctx, jobID, ai.JobStatusProcessing, ""); err != nil {
		workerLogger.Error("failed to update job status", "error", err)
		return
	}

	// 3. Get job details
	job, err := w.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		workerLogger.Error("failed to get job details", "error", err)
		w.handleJobFailure(ctx, jobID, "failed to get job details: "+err.Error())
		return
	}

	// 4. Run whisper transcription
	workerLogger.Debug("starting transcription", "audioPath", job.AudioPath)

	startTime := time.Now()
	result, err := w.whisperClient.Transcribe(ctx, job.AudioPath)
	processingTimeMs := time.Since(startTime).Milliseconds()

	if err != nil {
		workerLogger.Error("transcription failed", "error", err)
		w.handleJobFailure(ctx, jobID, "transcription failed: "+err.Error())
		return
	}

	workerLogger.Info("transcription completed",
		"duration", result.Duration,
		"confidence", result.ConfidenceScore,
		"processingTimeMs", processingTimeMs)

	// 5. Save transcription
	// Create domain object
	transcription := ai.NewTranscription(
		job.ID,
		job.MediaFileID,
		"", // AudioURL - would be set based on storage
		result.Transcript,
		result.ConfidenceScore,
		result.Language,
		result.Duration,
		"whisper.cpp", // ModelName
		processingTimeMs,
	)

	// Encrypt using encx
	transcriptionEncx, err := ai.ProcessTranscriptionEncx(ctx, w.crypto, transcription)
	if err != nil {
		workerLogger.Error("failed to encrypt transcription", "error", err)
		w.handleJobFailure(ctx, jobID, "failed to encrypt transcription: "+err.Error())
		return
	}

	// Persist encrypted transcription
	if err := w.transcriptionRepo.Create(ctx, transcriptionEncx); err != nil {
		workerLogger.Error("failed to save transcription", "error", err)
		w.handleJobFailure(ctx, jobID, "failed to save transcription: "+err.Error())
		return
	}

	// 6. Mark job complete
	if err := w.jobRepo.Complete(ctx, jobID, transcriptionEncx.ID); err != nil {
		workerLogger.Error("failed to mark job complete", "error", err)
		return
	}

	workerLogger.Info("job completed successfully", "transcriptionID", transcriptionEncx.ID)

	// 7. Send webhook if configured
	if job.WebhookURL != nil {
		if err := w.webhookClient.NotifyCompletion(ctx, *job.WebhookURL, job.ID, transcriptionEncx.ID); err != nil {
			workerLogger.Warn("failed to send webhook", "error", err)
			// Don't fail the job if webhook fails
		}
	}

	// 8. Cleanup audio file if configured
	if w.deleteAudioAfter {
		if err := os.Remove(job.AudioPath); err != nil && !os.IsNotExist(err) {
			workerLogger.Warn("failed to delete audio file", "error", err)
		} else {
			workerLogger.Debug("audio file deleted", "audioPath", job.AudioPath)
		}
	}
}

// handleJobFailure handles a job processing failure.
func (w *TranscriptionWorker) handleJobFailure(ctx context.Context, jobID uuid.UUID, errorMsg string) {
	if err := w.jobRepo.UpdateStatus(ctx, jobID, ai.JobStatusFailed, errorMsg); err != nil {
		w.logger.Error("failed to mark job as failed",
			"jobID", jobID,
			"error", err)
	}

	// Get job details for audio cleanup and webhook notification.
	job, err := w.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		w.logger.Error("failed to get job for failure cleanup", "jobID", jobID, "error", err)
		return
	}

	if err := os.Remove(job.AudioPath); err != nil && !os.IsNotExist(err) {
		w.logger.Warn("failed to delete audio file after job failure",
			"jobID", jobID,
			"audioPath", job.AudioPath,
			"error", err)
	}

	// Delete job before media: fk_jobs_media_file ON DELETE CASCADE would wipe the
	// job row if media were deleted first, defeating our explicit ordering.
	if err := w.jobRepo.Delete(ctx, jobID); err != nil {
		w.logger.Error("failed to delete failed job",
			"jobID", jobID,
			"error", err)
	}

	// Delete the media record so the recording does not appear stuck in
	// "transcribing" state on the home list.
	if err := w.mediaRepo.DeleteMedia(ctx, job.MediaFileID); err != nil {
		w.logger.Warn("failed to delete media after job failure",
			"jobID", jobID,
			"mediaFileID", job.MediaFileID,
			"error", err)
	}

	if job.WebhookURL != nil {
		if err := w.webhookClient.NotifyFailure(ctx, *job.WebhookURL, jobID, errorMsg); err != nil {
			w.logger.Warn("failed to send failure webhook",
				"jobID", jobID,
				"error", err)
		}
	}
}

// pollPendingJobs polls for pending jobs that haven't been enqueued.
// This handles job recovery on startup and picks up jobs that weren't enqueued via API.
func (w *TranscriptionWorker) pollPendingJobs(ctx context.Context) {
	defer w.wg.Done()

	w.logger.Info("starting pending jobs poller")

	// On startup, recover any pending/orphaned jobs
	w.recoverPendingJobs(ctx)

	// Periodic polling ticker
	ticker := time.NewTicker(DefaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			w.logger.Debug("pending jobs poller received stop signal")
			return
		case <-ctx.Done():
			w.logger.Debug("pending jobs poller context cancelled")
			return
		case <-ticker.C:
			// Poll for pending jobs
			jobs, err := w.jobRepo.GetPendingJobs(ctx, 10)
			if err != nil {
				w.logger.Error("failed to get pending jobs", "error", err)
				continue
			}

			// Enqueue any found jobs
			for _, job := range jobs {
				w.logger.Debug("enqueuing recovered job", "jobID", job.ID)
				w.Enqueue(job.ID)
			}
		}
	}
}

// recoverPendingJobs recovers pending jobs on startup.
func (w *TranscriptionWorker) recoverPendingJobs(ctx context.Context) {
	w.logger.Info("recovering pending jobs")

	jobs, err := w.jobRepo.GetPendingJobs(ctx, 100)
	if err != nil {
		w.logger.Error("failed to get pending jobs for recovery", "error", err)
		return
	}

	for _, job := range jobs {
		w.logger.Debug("enqueuing recovered job",
			"jobID", job.ID,
			"createdAt", job.CreatedAt)

		w.Enqueue(job.ID)
	}

	w.logger.Info("recovered pending jobs", "count", len(jobs))
}

// GetQueueSize returns the current queue size.
func (w *TranscriptionWorker) GetQueueSize() int {
	return w.queue.Size()
}

// GetWorkerID returns the worker ID.
func (w *TranscriptionWorker) GetWorkerID() string {
	return w.workerID
}
