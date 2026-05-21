package caseSubjectHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

func (h *handler) CreateCaseSubject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var req subject.CreateCaseSubjectRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		logger.WarnContext(ctx, "Handler: invalid request body", "operation", "create_case_subject", "error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.svc.CreateCaseSubject(ctx, &req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create case subject")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
