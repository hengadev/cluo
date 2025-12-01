package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

func (h *handler) GetContactByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract contact ID from URL path
	contactIDStr := r.PathValue("id")
	if contactIDStr == "" {
		logger.WarnContext(ctx, "Handler: contact ID is required",
			"operation", "get_contact_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
	}

	// Parse user ID as UUID
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid contact ID format",
			"operation", "get_contact_by_id",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Build the request
	request := &client.GetContactByIDRequest{
		ContactID: contactID,
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get contact request",
		"operation", "get_contact_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"contact_id", contactID,
		"user_agent", r.Header.Get("User-Agent"))

	contact, err := h.svc.GetContactByID(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get contact by ID")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get contact by ID request completed successfully",
		"operation", "get_contact_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with contact data
	httpx.RespondWithJSON(w, struct {
		Message string                  `json:"message"`
		Contact *client.ContactResponse `json:"contact"`
	}{
		Message: "Contact retrieval by ID completed successfully",
		Contact: contact,
	}, http.StatusOK)
}
