package rapportHandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (h *handler) UpdateRapport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawID := r.PathValue("id")
	caseID, err := uuid.Parse(rawID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, errs.NewInvalidValueErr("invalid case ID format"), "update rapport")
		return
	}

	var payload rapport.UpdateRapportRequest
	payload.CaseID = caseID

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_rapport",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing update rapport request",
		"operation", "update_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", rawID)

	response, err := h.svc.UpdateRapport(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update rapport")
		return
	}

	logger.InfoContext(ctx, "Handler: Update rapport request completed successfully",
		"operation", "update_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
