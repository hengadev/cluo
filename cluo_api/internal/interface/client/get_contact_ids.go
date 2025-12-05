package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"

	"github.com/google/uuid"
)

type GetContactIDsResponse struct {
	ClientID    string   `json:"client_id"`
	ContactIDs  []string `json:"contact_ids"`
	Count       int      `json:"count"`
}

func (h *handler) GetContactIDsForClient(w http.ResponseWriter, r *http.Request) {
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
			"operation", "get_contact_ids_for_client",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse client ID as UUID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid client ID format",
			"operation", "get_contact_ids_for_client",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get contact IDs request",
		"operation", "get_contact_ids_for_client",
		"method", r.Method,
		"path", r.URL.Path,
		"client_id", clientID,
		"user_agent", r.Header.Get("User-Agent"))

	contactIDs, err := h.svc.GetContactIDsForClient(ctx, clientID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get contact IDs for client")
		return
	}

	// Convert UUIDs to strings for JSON response
	contactIDStrings := make([]string, len(contactIDs))
	for i, id := range contactIDs {
		contactIDStrings[i] = id.String()
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get contact IDs request completed successfully",
		"operation", "get_contact_ids_for_client",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"contact_count", len(contactIDs))

	// Respond with contact IDs
	response := GetContactIDsResponse{
		ClientID:   clientID.String(),
		ContactIDs: contactIDStrings,
		Count:      len(contactIDs),
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}