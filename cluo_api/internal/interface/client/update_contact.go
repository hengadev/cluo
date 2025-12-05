package clientHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

func (h *handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract client ID from URL path
	contactIDStr := r.PathValue("id")
	if contactIDStr == "" {
		logger.WarnContext(ctx, "Handler: client ID is required",
			"operation", "update_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
	}

	// Parse user ID as UUID
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid contact ID format",
			"operation", "update_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	var payload client.UpdateContactRequest
	payload.ID = contactID

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_contact",
			"method", r.Method,
			"path", r.URL.Path)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing update contact request",
		"operation", "update_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"contact_id", contactID,
		"user_agent", r.Header.Get("User-Agent"))

	res, err := h.svc.UpdateContact(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update contact")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Update contact request completed successfully",
		"operation", "update_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, res, http.StatusOK)
}
