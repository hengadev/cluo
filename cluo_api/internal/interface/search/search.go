package searchHandler

import (
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("q")
	if len(query) > 1000 {
		httpx.RespondWithError(w, fmt.Errorf("query too long"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing search request",
		"operation", "search",
		"method", r.Method,
		"path", r.URL.Path,
		"query_length", len(query))

	resp, err := h.svc.Search(ctx, query)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "search")
		return
	}

	logger.InfoContext(ctx, "Handler: Search completed",
		"operation", "search",
		"results_count", len(resp.Results))

	httpx.RespondWithJSON(w, resp, http.StatusOK)
}
