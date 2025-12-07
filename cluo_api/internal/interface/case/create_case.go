package caseHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (h *handler) CreateCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload caseDomain.CreateCaseRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_case",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create case request",
		"operation", "create_case",
		"method", r.Method,
		"path", r.URL.Path,
		"case_title", payload.Title,
		"client_id", payload.ClientID,
		"user_agent", r.Header.Get("User-Agent"))

	response, err := h.svc.CreateCase(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create case")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Create case request completed successfully",
		"operation", "create_case",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
