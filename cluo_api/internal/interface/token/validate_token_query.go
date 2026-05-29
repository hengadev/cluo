package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

// ValidateTokenQuery handles GET /tokens/validate?token=<raw-token>.
// It validates the token and returns the associated case ID, or an error if
// the token is invalid or expired.
func (h *TokenHandler) ValidateTokenQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.URL.Query().Get("token")
	if rawToken == "" {
		httpx.RespondWithError(w, errMissingTokenParam, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing validate token query request",
		"operation", "validate_token_query",
		"method", r.Method,
		"path", r.URL.Path)

	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate token query")
		return
	}

	httpx.RespondWithJSON(w, validateTokenResponse{CaseID: caseID.String()}, http.StatusOK)
}
