package tokenHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) GetMediaByIDByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")
	rawMediaID := r.PathValue("mediaId")

	mediaID, err := uuid.Parse(rawMediaID)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("mediaId must be a valid UUID"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get media by id by token request",
		"operation", "get_media_by_id_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	response, err := h.svc.GetPublishedMediaByIDAndToken(ctx, rawToken, mediaID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get media by id by token")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
