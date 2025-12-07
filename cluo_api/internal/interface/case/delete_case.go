package caseHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (h *handler) DeleteCase(w http.ResponseWriter, r *http.Request) {
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
			"operation", "delete_case",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse case ID as UUID
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid case ID format",
			"operation", "delete_case",
			"method", r.Method,
			"path", r.URL.Path,
			"case_id", caseIDStr,
			"error", err)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete case request",
		"operation", "delete_case",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID,
		"user_agent", r.Header.Get("User-Agent"))

	// Create request object
	request := caseDomain.DeleteCaseByIDRequest{
		ID: caseID,
	}

	// Call service layer
	err = h.svc.DeleteCase(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete case")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Delete case request completed successfully",
		"operation", "delete_case",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusNoContent,
		"case_id", caseID)

	// Return success with no content
	w.WriteHeader(http.StatusNoContent)
}