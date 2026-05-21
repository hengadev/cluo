package tokenHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) ListTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	caseIDStr := r.PathValue("id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid case ID",
			"error", err,
			"operation", "list_tokens",
			"case_id", caseIDStr)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing list tokens request",
		"operation", "list_tokens",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseIDStr)

	response, err := h.svc.ListTokensByCaseID(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list tokens")
		return
	}

	logger.InfoContext(ctx, "Handler: List tokens request completed successfully",
		"operation", "list_tokens",
		"case_id", caseIDStr,
		"count", len(response))

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
