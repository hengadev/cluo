package caseHandler

import (
	"net/http"
	"strconv"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (h *handler) ListCases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Only accept GET method
	if r.Method != http.MethodGet {
		logger.WarnContext(ctx, "Handler: Method not allowed",
			"operation", "list_cases",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, nil, http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for filtering and pagination
	query := r.URL.Query()

	// Parse pagination parameters with validation and defaults
	page, err := parsePositiveIntQueryParam(query.Get("page"), 1, "page")
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid page parameter",
			"operation", "list_cases",
			"error", err,
			"page_param", query.Get("page"))
		httpx.RespondWithError(w, nil, http.StatusBadRequest)
		return
	}

	pageSize, err := parsePositiveIntQueryParam(query.Get("pageSize"), 20, "pageSize")
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid pageSize parameter",
			"operation", "list_cases",
			"error", err,
			"page_size_param", query.Get("pageSize"))
		httpx.RespondWithError(w, nil, http.StatusBadRequest)
		return
	}

	// Validate page size limits
	if pageSize > 100 {
		logger.WarnContext(ctx, "Handler: Page size too large",
			"operation", "list_cases",
			"page_size", pageSize)
		httpx.RespondWithError(w, nil, http.StatusBadRequest)
		return
	}

	// Create request object
	var request caseDomain.ListCasesRequest
	request.Page = page
	request.PageSize = pageSize

	// Parse optional filter parameters
	if clientID := query.Get("clientId"); clientID != "" {
		request.ClientID = &clientID
	}
	if status := query.Get("status"); status != "" {
		request.Status = &status
	}
	if assignedContactID := query.Get("assignedContactId"); assignedContactID != "" {
		request.AssignedContactID = &assignedContactID
	}
	if dateCreatedFrom := query.Get("dateCreatedFrom"); dateCreatedFrom != "" {
		request.DateCreatedFrom = &dateCreatedFrom
	}
	if dateCreatedTo := query.Get("dateCreatedTo"); dateCreatedTo != "" {
		request.DateCreatedTo = &dateCreatedTo
	}
	if dateUpdatedFrom := query.Get("dateUpdatedFrom"); dateUpdatedFrom != "" {
		request.DateUpdatedFrom = &dateUpdatedFrom
	}
	if dateUpdatedTo := query.Get("dateUpdatedTo"); dateUpdatedTo != "" {
		request.DateUpdatedTo = &dateUpdatedTo
	}
	if search := query.Get("search"); search != "" {
		// Validate search length
		if len(search) > 1000 {
			logger.WarnContext(ctx, "Handler: Search term too long",
				"operation", "list_cases",
				"search_length", len(search))
			httpx.RespondWithError(w, nil, http.StatusBadRequest)
			return
		}
		request.Search = &search
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing list cases request",
		"operation", "list_cases",
		"method", r.Method,
		"path", r.URL.Path,
		"page", page,
		"page_size", pageSize,
		"has_client_filter", request.ClientID != nil,
		"has_status_filter", request.Status != nil,
		"has_search", request.Search != nil,
		"user_agent", r.Header.Get("User-Agent"))

	// Validate request using domain validation
	if err := request.Valid(ctx); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid request parameters",
			"operation", "list_cases",
			"error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Call service layer
	response, err := h.svc.List(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list cases")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: List cases request completed successfully",
		"operation", "list_cases",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"total_cases", response.Pagination.TotalItems,
		"returned_cases", len(response.Cases),
		"page", response.Pagination.Page,
		"total_pages", response.Pagination.TotalPages)

	// Return response
	httpx.RespondWithJSON(w, response, http.StatusOK)
}

// parsePositiveIntQueryParam parses a query parameter as a positive integer with validation
func parsePositiveIntQueryParam(value string, defaultValue int, paramName string) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, &InvalidQueryParamError{
			Parameter: paramName,
			Value:     value,
			Message:   "must be a valid integer",
		}
	}

	if parsed < 1 {
		return defaultValue, &InvalidQueryParamError{
			Parameter: paramName,
			Value:     value,
			Message:   "must be greater than 0",
		}
	}

	return parsed, nil
}

// InvalidQueryParamError represents an error with a query parameter
type InvalidQueryParamError struct {
	Parameter string
	Value     string
	Message   string
}

func (e *InvalidQueryParamError) Error() string {
	return e.Parameter + ": " + e.Message + " (got: " + e.Value + ")"
}