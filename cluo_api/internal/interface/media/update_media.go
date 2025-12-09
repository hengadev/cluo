package mediaHandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (h *handler) UpdateMedia(w http.ResponseWriter, r *http.Request) {
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
			"operation", "update_media",
			"id", idStr,
			"error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	var payload domain.UpdateMediaRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "update_media",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	payload.ID = mediaID

	logger.InfoContext(ctx, "Handler: Processing update media request",
		"operation", "update_media",
		"media_id", mediaID.String())

	response, err := h.svc.UpdateMedia(ctx, &payload)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update media")
		return
	}

	logger.InfoContext(ctx, "Handler: Update media completed successfully",
		"operation", "update_media",
		"media_id", response.ID)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
