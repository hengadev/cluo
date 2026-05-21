package pieceHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) GetPiece(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	pieceIDStr := r.PathValue("pieceId")
	pieceID, err := uuid.Parse(pieceIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid piece ID",
			"operation", "get_piece",
			"pieceId", pieceIDStr,
			"error", err)
		httpx.RespondWithError(w, fmt.Errorf("invalid piece ID"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get piece request",
		"operation", "get_piece",
		"piece_id", pieceID.String())

	response, err := h.svc.GetPieceByID(ctx, pieceID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get piece")
		return
	}

	logger.InfoContext(ctx, "Handler: Get piece completed successfully",
		"operation", "get_piece",
		"piece_id", response.ID)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
