package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"

	"github.com/google/uuid"
)

func (h *handler) GetAllContactsByClientID(w http.ResponseWriter, r *http.Request) {

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
			"operation", "get_all_contacts_by_contact_id", "method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse user ID as UUID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid client ID format",
			"operation", "get_all_contacts_by_contact_id", "method", r.Method,
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	contacts, err := h.svc.GetAllContactsByClientID(ctx, clientID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all contacts by client ID")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get all contacts by client ID request completed successfully",
		"operation", "get_all_contacts_by_client_id",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, contacts, http.StatusOK)
}
