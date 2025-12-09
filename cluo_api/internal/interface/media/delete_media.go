package mediaHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (h *handler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
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
			"operation", "delete_media",
			"id", idStr,
			"error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	request := &domain.DeleteMediaRequest{ID: mediaID}

	logger.InfoContext(ctx, "Handler: Processing delete media request",
		"operation", "delete_media",
		"media_id", mediaID.String())

	if err := h.svc.DeleteMedia(ctx, request); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete media")
		return
	}

	logger.InfoContext(ctx, "Handler: Delete media completed successfully",
		"operation", "delete_media",
		"media_id", mediaID.String())

	w.WriteHeader(http.StatusNoContent)
}
