package rapportHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) GetRapport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawID := r.PathValue("id")
	caseID, err := uuid.Parse(rawID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, errs.NewInvalidValueErr("invalid case ID format"), "get rapport")
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get rapport request",
		"operation", "get_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", rawID)

	response, err := h.svc.GetRapportByCaseID(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get rapport")
		return
	}

	logger.InfoContext(ctx, "Handler: Get rapport request completed successfully",
		"operation", "get_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
