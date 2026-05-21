package rapportHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) DeleteRapport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawID := r.PathValue("id")
	caseID, err := uuid.Parse(rawID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, errs.NewInvalidValueErr("invalid case ID format"), "delete rapport")
		return
	}

	logger.InfoContext(ctx, "Handler: Processing delete rapport request",
		"operation", "delete_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", rawID)

	if err := h.svc.DeleteRapport(ctx, caseID); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete rapport")
		return
	}

	logger.InfoContext(ctx, "Handler: Delete rapport request completed successfully",
		"operation", "delete_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusNoContent)

	w.WriteHeader(http.StatusNoContent)
}
