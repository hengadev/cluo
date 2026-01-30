package aiSpeechToTextHandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
)

type ListJobsResponse struct {
	Jobs  []JobSummary `json:"jobs"`
	Total int          `json:"total"`
	Limit int          `json:"limit"`
	Offset int         `json:"offset"`
}

type JobSummary struct {
	JobID       uuid.UUID    `json:"jobId"`
	MediaFileID uuid.UUID    `json:"mediaFileId"`
	Status      ai.JobStatus `json:"status"`
	Progress    int          `json:"progress"`
	CreatedAt   string       `json:"createdAt"`
}

func (h *handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get user ID from context
	userID, err := ctxutil.GetUserIDFromContext(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: User ID not found in context",
			"error", err,
			"operation", "list_jobs")
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	statusStr := r.URL.Query().Get("status")

	limit := 20
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
		if limit > 100 {
			limit = 100
		}
	}

	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	var status *ai.JobStatus
	if statusStr != "" {
		s := ai.JobStatus(statusStr)
		if s.Valid() {
			status = &s
		}
	}

	logger.DebugContext(ctx, "Handler: Processing list jobs request",
		"operation", "list_jobs",
		"userID", userID,
		"limit", limit,
		"offset", offset,
		"status", statusStr)

	// Call service
	req := &ports.ListJobsRequest{
		UserID:  userID,
		Status:  status,
		Limit:   limit,
		Offset:  offset,
	}

	jobs, total, err := h.svc.ListJobs(ctx, req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list jobs")
		return
	}

	// Build response
	jobSummaries := make([]JobSummary, len(jobs))
	for i, job := range jobs {
		jobSummaries[i] = JobSummary{
			JobID:       job.ID,
			MediaFileID: job.MediaFileID,
			Status:      job.Status,
			Progress:    job.Progress,
			CreatedAt:   job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListJobsResponse{
		Jobs:   jobSummaries,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	// Set JSON encoder to escape HTML for security
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	encoder.Encode(response)
}
