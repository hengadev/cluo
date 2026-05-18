package investigationHandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (h *handler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract case ID from URL path
	caseIDStr := r.PathValue("id")
	if caseIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing case ID in path",
			"operation", "update_case",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse case ID as UUID
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid case ID format",
			"operation", "update_case",
			"method", r.Method,
			"path", r.URL.Path,
			"case_id", caseIDStr,
			"error", err)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse request body
	var request investigation.UpdateCaseRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"operation", "update_case",
			"method", r.Method,
			"path", r.URL.Path,
			"error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Set the case ID from URL path
	request.ID = caseID

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing update case request",
		"operation", "update_case",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID,
		"user_agent", r.Header.Get("User-Agent"))

	// Call service layer
	response, err := h.svc.UpdateCase(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update case")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Update case request completed successfully",
		"operation", "update_case",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"case_id", caseID)

	// Return response
	httpx.RespondWithJSON(w, response, http.StatusOK)
}