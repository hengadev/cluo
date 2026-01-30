package aiSpeechToTextHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) CancelJob(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing cancel job request",
		"operation", "cancel_job",
		"jobID", jobID)

	// Call service
	if err := h.svc.CancelJob(ctx, jobID); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "cancel job")
		return
	}

	logger.InfoContext(ctx, "Handler: Cancel job request completed successfully",
		"operation", "cancel_job",
		"jobID", jobID)

	httpx.RespondWithJSON(w, map[string]string{"message": "Job cancelled"}, http.StatusOK)
}
