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

func (h *handler) CreateContact(w http.ResponseWriter, r *http.Request) {
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
			"operation", "create_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	// Parse user ID as UUID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid client ID format",
			"operation", "create_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
		return
	}

	var payload client.CreateContactRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	payload.ClientID = clientID

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create contact request",
		"operation", "create_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"client_id", clientID,
		"user_agent", r.Header.Get("User-Agent"))

	res, err := h.svc.CreateContact(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create contact")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Create contact request completed successfully",
		"operation", "create_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, res, http.StatusOK)
}
