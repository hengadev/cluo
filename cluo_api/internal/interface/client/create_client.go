package clientHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload client.CreateClientRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_client",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create client request",
		"operation", "create_client",
		"method", r.Method,
		"path", r.URL.Path,
		"client_name", payload.Name,
		"client_type", payload.Type,
		"user_agent", r.Header.Get("User-Agent"))

	if err = h.svc.CreateClient(ctx, &payload); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create client")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Create client request completed successfully",
		"operation", "create_client",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, map[string]string{"message": "Client creation completed successfully"}, http.StatusOK)
}

