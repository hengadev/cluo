package caseSubjectHandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

func (h *handler) UpdateCaseSubject(w http.ResponseWriter, r *http.Request) {
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

	var req subject.UpdateCaseSubjectRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		logger.WarnContext(ctx, "Handler: invalid request body", "operation", "update_case_subject", "error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}
	req.ID = id

	response, err := h.svc.UpdateCaseSubject(ctx, &req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update case subject")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
