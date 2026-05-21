package caseTypeHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (h *handler) CreateCaseType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var req casetype.CreateCaseTypeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		logger.WarnContext(ctx, "Handler: invalid request body", "operation", "create_case_type", "error", err)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.svc.CreateCaseType(ctx, &req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create case type")
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
