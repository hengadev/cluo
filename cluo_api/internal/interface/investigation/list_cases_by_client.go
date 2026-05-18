package investigationHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (h *handler) ListCasesByClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Only accept GET method
	if r.Method != http.MethodGet {
		logger.WarnContext(ctx, "Handler: Method not allowed",
			"operation", "list_cases_by_client",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// Extract client ID from URL path
	clientIDStr := r.PathValue("clientId")
	if clientIDStr == "" {
		logger.WarnContext(ctx, "Handler: Missing client ID in path",
			"operation", "list_cases_by_client",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, fmt.Errorf("client ID is required"), http.StatusBadRequest)
		return
	}

	// Validate client ID as UUID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid client ID format",
			"operation", "list_cases_by_client",
			"method", r.Method,
			"path", r.URL.Path,
			"client_id", clientIDStr,
			"error", err)
		httpx.RespondWithError(w, fmt.Errorf("invalid client ID format"), http.StatusBadRequest)
		return
	}

	// Parse query parameters for pagination
	query := r.URL.Query()

	// Parse pagination parameters with validation and defaults
	page, err := parsePositiveIntQueryParam(query.Get("page"), 1, "page")
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid page parameter",
			"operation", "list_cases_by_client",
			"error", err,
			"page_param", query.Get("page"))
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	pageSize, err := parsePositiveIntQueryParam(query.Get("pageSize"), 20, "pageSize")
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid pageSize parameter",
			"operation", "list_cases_by_client",
			"error", err,
			"page_size_param", query.Get("pageSize"))
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Validate page size limits
	if pageSize > 100 {
		logger.WarnContext(ctx, "Handler: Page size too large",
			"operation", "list_cases_by_client",
			"page_size", pageSize)
		httpx.RespondWithError(w, fmt.Errorf("page size cannot exceed 100"), http.StatusBadRequest)
		return
	}

	// Create request object
	request := investigation.ListByClientRequest{
		ClientID: clientIDStr,
		Page:     page,
		PageSize: pageSize,
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing list cases by client request",
		"operation", "list_cases_by_client",
		"method", r.Method,
		"path", r.URL.Path,
		"client_id", clientID,
		"page", page,
		"page_size", pageSize,
		"user_agent", r.Header.Get("User-Agent"))

	// Validate request using domain validation
	if err := request.Valid(ctx); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid request parameters",
			"operation", "list_cases_by_client",
			"error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Call service layer
	response, err := h.svc.ListByClient(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list cases by client")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: List cases by client request completed successfully",
		"operation", "list_cases_by_client",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"client_id", clientID,
		"total_cases", response.Pagination.TotalItems,
		"returned_cases", len(response.Cases),
		"page", response.Pagination.Page,
		"total_pages", response.Pagination.TotalPages)

	// Return response
	httpx.RespondWithJSON(w, response, http.StatusOK)
}