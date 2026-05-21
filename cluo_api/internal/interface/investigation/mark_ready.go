package investigationHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) MarkReady(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	caseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.RespondWithError(w, fmt.Errorf("invalid case ID"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing mark-ready request",
		"operation", "mark_ready",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID)

	response, err := h.svc.MarkReady(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "mark case ready")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
