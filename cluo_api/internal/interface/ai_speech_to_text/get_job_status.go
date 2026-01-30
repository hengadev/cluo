package aiSpeechToTextHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

type JobStatusResponse struct {
	JobID           uuid.UUID       `json:"jobId"`
	MediaFileID     uuid.UUID       `json:"mediaFileId"`
	Status          ai.JobStatus    `json:"status"`
	Progress        int             `json:"progress"`
	ErrorMessage    string          `json:"errorMessage,omitempty"`
	TranscriptionID *uuid.UUID      `json:"transcriptionId,omitempty"`
	CreatedAt       string          `json:"createdAt"`
	StartedAt       *string         `json:"startedAt,omitempty"`
	CompletedAt     *string         `json:"completedAt,omitempty"`
}

func (h *handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract job ID from URL path
	jobIDStr := r.PathValue("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.DebugContext(ctx, "Handler: Processing get job status request",
		"operation", "get_job_status",
		"jobID", jobID)

	// Call service
	job, err := h.svc.GetJobStatus(ctx, jobID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get job status")
		return
	}

	response := &JobStatusResponse{
		JobID:        job.ID,
		MediaFileID:  job.MediaFileID,
		Status:       job.Status,
		Progress:     job.Progress,
		ErrorMessage: job.ErrorMessage,
		TranscriptionID: job.TranscriptionID,
		CreatedAt:    job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if job.StartedAt != nil {
		startedAt := job.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		response.StartedAt = &startedAt
	}

	if job.CompletedAt != nil {
		completedAt := job.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAt
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
