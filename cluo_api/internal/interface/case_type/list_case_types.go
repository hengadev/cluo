package caseTypeHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) ListCaseTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	response, err := h.svc.ListCaseTypes(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list case types")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
