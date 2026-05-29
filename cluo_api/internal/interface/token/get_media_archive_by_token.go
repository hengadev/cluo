package tokenHandler

import (
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/archive"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *TokenHandler) GetMediaArchiveByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")

	logger.InfoContext(ctx, "Handler: Processing get media archive by token request",
		"operation", "get_media_archive_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	// Validate token and retrieve caseID
	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate token for media archive")
		return
	}

	if h.archiveDeps == nil {
		httpx.RespondWithError(w, fmt.Errorf("archive feature not available"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="medias.zip"`)

	if err := archive.BuildMediaArchive(ctx, h.archiveDeps, caseID, w); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to build media archive",
			"error", err.Error(),
			"operation", "get_media_archive_by_token",
		)
		return
	}

	logger.InfoContext(ctx, "Handler: Get media archive by token completed successfully",
		"operation", "get_media_archive_by_token",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)
}
