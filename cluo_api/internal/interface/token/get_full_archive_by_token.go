package tokenHandler

import (
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/archive"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (h *TokenHandler) GetFullArchiveByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")

	logger.InfoContext(ctx, "Handler: Processing get full archive by token request",
		"operation", "get_full_archive_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	// Validate token and retrieve caseID
	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate token for full archive")
		return
	}

	if h.archiveDeps == nil {
		httpx.RespondWithError(w, fmt.Errorf("archive feature not available"), http.StatusNotFound)
		return
	}

	// Use the case external reference for a human-readable filename; fall back to the UUID.
	fileName := fmt.Sprintf("dossier-%s.zip", caseID)
	if h.caseSvc != nil {
		if caseResp, err := h.caseSvc.GetCaseByID(ctx, &investigation.GetCaseByIDRequest{ID: caseID}); err == nil &&
			caseResp.ExternalReference != nil && *caseResp.ExternalReference != "" {
			fileName = fmt.Sprintf("dossier-%s.zip", *caseResp.ExternalReference)
		}
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

	if err := archive.BuildFullArchive(ctx, h.archiveDeps, caseID, w); err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to build full archive",
			"error", err.Error(),
			"operation", "get_full_archive_by_token",
		)
		// Headers are already sent — status code cannot be changed; error is logged for monitoring.
		return
	}

	logger.InfoContext(ctx, "Handler: Get full archive by token completed successfully",
		"operation", "get_full_archive_by_token",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)
}
