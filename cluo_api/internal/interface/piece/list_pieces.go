package pieceHandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) ListPieces(w http.ResponseWriter, r *http.Request) {
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
			"operation", "list_pieces",
			"id", caseIDStr,
			"error", err)
		httpx.RespondWithError(w, fmt.Errorf("invalid case ID"), http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logger.InfoContext(ctx, "Handler: Processing list pieces request",
		"operation", "list_pieces",
		"case_id", caseID.String(),
		"page", page,
		"page_size", pageSize)

	response, err := h.svc.ListPiecesByCaseID(ctx, caseID, page, pageSize)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list pieces")
		return
	}

	logger.InfoContext(ctx, "Handler: List pieces completed successfully",
		"operation", "list_pieces",
		"case_id", caseID.String(),
		"count", len(response.Pieces))

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
