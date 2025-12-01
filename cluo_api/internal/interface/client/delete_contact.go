package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

func (h *handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
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
			"operation", "delete_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse user ID as UUID
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid contact ID format",
			"operation", "delete_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing delete contact request",
		"operation", "delete_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"contact_id", contactID,
		"user_agent", r.Header.Get("User-Agent"))

	payload := &client.DeleteContactRequest{
		ContactID: contactID,
	}

	if err = h.svc.DeleteContact(ctx, payload); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete contact")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Delete contact request completed successfully",
		"operation", "delete_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "Contact deletion completed successfully",
	}, http.StatusOK)
}
