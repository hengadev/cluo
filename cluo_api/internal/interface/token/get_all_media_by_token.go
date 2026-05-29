package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *TokenHandler) GetAllMediaByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")

	logger.InfoContext(ctx, "Handler: Processing get all media by token request",
		"operation", "get_all_media_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetPublishedMediaByToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get all media by token")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
