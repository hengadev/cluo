package pieceHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) DeletePiece(w http.ResponseWriter, r *http.Request) {
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
			"operation", "delete_piece",
			"id", caseIDStr,
			"error", err)
		httpx.RespondWithError(w, fmt.Errorf("invalid case ID"), http.StatusBadRequest)
		return
	}

	pieceIDStr := r.PathValue("pieceId")
	pieceID, err := uuid.Parse(pieceIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid piece ID",
			"operation", "delete_piece",
			"pieceId", pieceIDStr,
			"error", err)
		httpx.RespondWithError(w, fmt.Errorf("invalid piece ID"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing delete piece request",
		"operation", "delete_piece",
		"case_id", caseID.String(),
		"piece_id", pieceID.String())

	if err := h.svc.DeletePiece(ctx, caseID, pieceID); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete piece")
		return
	}

	logger.InfoContext(ctx, "Handler: Delete piece completed successfully",
		"operation", "delete_piece",
		"piece_id", pieceID.String())

	w.WriteHeader(http.StatusNoContent)
}
