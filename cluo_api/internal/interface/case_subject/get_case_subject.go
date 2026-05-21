package caseSubjectHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) GetCaseSubject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.RespondWithError(w, fmt.Errorf("invalid subject ID"), http.StatusBadRequest)
		return
	}

	response, err := h.svc.GetCaseSubjectByID(ctx, id)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get case subject")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
