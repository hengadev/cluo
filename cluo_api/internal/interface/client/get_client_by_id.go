package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

func (h *handler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract client ID from URL path
	clientIDStr := r.PathValue("id")
	if clientIDStr == "" {
		logger.WarnContext(ctx, "Handler: client ID is required",
			"operation", "get_client_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse client ID as UUID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid client ID format",
			"operation", "get_client_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Build the request
	request := &client.GetClientByIDRequest{
		ID: clientID,
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get client request",
		"operation", "get_client_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"client_id", clientID,
		"user_agent", r.Header.Get("User-Agent"))

	res, err := h.svc.GetClientByID(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get client by ID")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get client by ID request completed successfully",
		"operation", "get_client_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with client data
	httpx.RespondWithJSON(w, struct {
		Message string                 `json:"message"`
		Client  *client.ClientResponse `json:"client"`
	}{
		Message: "Client retrieval by ID completed successfully",
		Client:  res,
	}, http.StatusOK)
}

