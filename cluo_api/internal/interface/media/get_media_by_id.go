package mediaHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (h *handler) GetMediaByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse ID from path
	idStr := r.PathValue("id")
	mediaID, err := uuid.Parse(idStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid media ID",
			"operation", "get_media_by_id",
			"id", idStr,
			"error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	request := &domain.GetMediaByIDRequest{ID: mediaID}

	logger.InfoContext(ctx, "Handler: Processing get media by ID request",
		"operation", "get_media_by_id",
		"media_id", mediaID.String())

	response, err := h.svc.GetMediaByID(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get media")
		return
	}

	logger.InfoContext(ctx, "Handler: Get media by ID completed successfully",
		"operation", "get_media_by_id",
		"media_id", response.ID)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
