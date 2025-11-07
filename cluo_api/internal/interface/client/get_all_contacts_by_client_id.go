package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (h *handler) GetAllContactsByClientID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract client ID from URL path
	clientID := r.PathValue("id")
	if clientID == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("client ID is required"), http.StatusBadRequest)
		return
	}

	contacts, err := h.svc.GetAllContactsByClientID(ctx, clientID)
	if err != nil {

		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		}
		logFields := []any{
			"operation", "get_all_contacts_by_client_id",
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
			logger.InfoContext(ctx, "Handler: Get all contacts by client ID request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Get all contacts by client ID request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Get all contacts by client ID request failed", logFields...)
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get all contacts by client ID request completed successfully",
		"operation", "get_all_contacts_by_client_id",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, struct {
		Message  string                    `json:"message"`
		Contacts []*client.ContactResponse `json:"contacts"`
	}{
		Message:  "Contacts retrieval by client ID completed successfully",
		Contacts: contacts,
	}, http.StatusOK)
}
