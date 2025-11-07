package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (h *handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract client ID from URL path
	contactID := r.PathValue("id")
	if contactID == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("contact ID is required"), http.StatusBadRequest)
		return
	}

	var payload client.UpdateContactRequest
	payload.ID = contactID

	// TODO: need to decode the payload to get the rest of the request brother

	if err = h.svc.UpdateContact(ctx, &payload); err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		}
		logFields := []any{
			"operation", "update_contact",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}
		var statusCode int
		switch {
		}

		logFields = append(logFields, "status_code", statusCode)

		switch logLevel {
		case "info":
			logger.InfoContext(ctx, "Handler: update contact request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: update contact request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: update contact request failed", logFields...)
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Update contact request completed successfully",
		"operation", "update_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "Contact update completed successfully",
	}, http.StatusOK)
}
