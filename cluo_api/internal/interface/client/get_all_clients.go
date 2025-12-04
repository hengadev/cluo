package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (h *handler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get all clients request",
		"operation", "get_all_clients",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	clients, err := h.svc.GetAllClients(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all clients")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get all clients request completed successfully",
		"operation", "get_all_clients",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"clients_count", len(clients))

	// Respond with success message and clients data
	httpx.RespondWithJSON(w, struct {
		Message string                   `json:"message"`
		Clients []*client.ClientResponse `json:"clients"`
	}{
		Message: "All clients retrieval completed successfully",
		Clients: clients,
	}, http.StatusOK)
}

