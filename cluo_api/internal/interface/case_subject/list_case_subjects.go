package caseSubjectHandler

import (
	"net/http"
	"strconv"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) ListCaseSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	page := 1
	pageSize := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	responses, total, err := h.svc.ListCaseSubjects(ctx, page, pageSize)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list case subjects")
		return
	}

	httpx.RespondWithJSON(w, map[string]interface{}{
		"data":     responses,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, http.StatusOK)
}
