package caseTypeHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) DeleteCaseType(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.DeleteCaseType(ctx, id); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete case type")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
