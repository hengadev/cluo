package aiTranscriptAnalysisHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Check if service is available
	if h.svc == nil {
		logger.WarnContext(ctx, "Handler: Transcript analysis service unavailable")

		response := HealthCheckResponse{
			Status:  "unhealthy",
			Message: "service not initialized",
		}
		httpx.RespondWithJSON(w, response, http.StatusServiceUnavailable)
		return
	}

	response := HealthCheckResponse{
		Status: "healthy",
	}
	httpx.RespondWithJSON(w, response, http.StatusOK)
}
