package tokenHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	tokenIDStr := r.PathValue("tokenId")
	tokenID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid token ID",
			"error", err,
			"operation", "revoke_token",
			"token_id", tokenIDStr)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing revoke token request",
		"operation", "revoke_token",
		"method", r.Method,
		"path", r.URL.Path,
		"token_id", tokenIDStr)

	if err := h.svc.RevokeToken(ctx, tokenID); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "revoke token")
		return
	}

	logger.InfoContext(ctx, "Handler: Revoke token request completed successfully",
		"operation", "revoke_token",
		"token_id", tokenIDStr)

	w.WriteHeader(http.StatusNoContent)
}
