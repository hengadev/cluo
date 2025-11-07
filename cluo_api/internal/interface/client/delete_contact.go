package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

// TODO: fill the part with the error handling on the service call.

func (h *handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract contact ID from URL path
	contactID := r.PathValue("id")
	if contactID == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("contact ID is required"), http.StatusBadRequest)
		return
	}

	payload := &client.DeleteContactRequest{
		ContactID: contactID,
	}

	if err = h.svc.DeleteContact(ctx, payload); err != nil {

		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		}

		logFields := []any{
			"operation", "delete_contact",
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
			logger.InfoContext(ctx, "Handler: delete contact request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: delete contact request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: delete contact request failed", logFields...)
		}
		httpx.RespondWithError(w, err, statusCode)
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
