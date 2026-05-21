package rapportHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (h *handler) CreateRapport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	caseID := r.PathValue("id")

	var payload rapport.CreateRapportRequest
	payload.CaseID = caseID

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_rapport",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing create rapport request",
		"operation", "create_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID)

	response, err := h.svc.CreateRapport(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create rapport")
		return
	}

	logger.InfoContext(ctx, "Handler: Create rapport request completed successfully",
		"operation", "create_rapport",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
