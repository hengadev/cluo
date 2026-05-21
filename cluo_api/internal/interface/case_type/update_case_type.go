package caseTypeHandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (h *handler) UpdateCaseType(w http.ResponseWriter, r *http.Request) {
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

	var req casetype.UpdateCaseTypeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		logger.WarnContext(ctx, "Handler: invalid request body", "operation", "update_case_type", "error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.svc.UpdateCaseType(ctx, id, &req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update case type")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
