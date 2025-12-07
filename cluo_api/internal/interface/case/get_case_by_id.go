package caseHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (h *handler) GetCaseByID(w http.ResponseWriter, r *http.Request) {
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
			"operation", "get_case_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse case ID as UUID
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid case ID format",
			"operation", "get_case_by_id",
			"method", r.Method,
			"path", r.URL.Path,
			"case_id", caseIDStr,
			"error", err)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get case by ID request",
		"operation", "get_case_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID,
		"user_agent", r.Header.Get("User-Agent"))

	// Create request object
	request := caseDomain.GetCaseByIDRequest{
		ID: caseID,
	}

	// Call service layer
	response, err := h.svc.GetCaseByID(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get case by ID")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get case by ID request completed successfully",
		"operation", "get_case_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"case_id", caseID)

	// Return response
	httpx.RespondWithJSON(w, response, http.StatusOK)
}

