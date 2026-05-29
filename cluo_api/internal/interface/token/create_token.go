package tokenHandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *TokenHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	caseIDStr := r.PathValue("id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid case ID",
			"error", err,
			"operation", "create_token",
			"case_id", caseIDStr)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing create token request",
		"operation", "create_token",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseIDStr)

	response, err := h.svc.CreateToken(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create token")
		return
	}

	logger.InfoContext(ctx, "Handler: Create token request completed successfully",
		"operation", "create_token",
		"case_id", caseIDStr,
		"token_id", response.ID)

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
