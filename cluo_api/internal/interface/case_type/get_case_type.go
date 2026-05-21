package caseTypeHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) GetCaseType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.RespondWithError(w, fmt.Errorf("invalid case type ID"), http.StatusBadRequest)
		return
	}

	response, err := h.svc.GetCaseTypeByID(ctx, id)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get case type")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
