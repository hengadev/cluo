package pieceHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

const maxUploadSize = 100 << 20 // 100 MB

func (h *handler) UploadPiece(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse case ID from path
	caseIDStr := r.PathValue("id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid case ID",
			"operation", "upload_piece",
			"id", caseIDStr,
			"error", err)
		httpx.RespondWithError(w, fmt.Errorf("invalid case ID"), http.StatusBadRequest)
		return
	}

	// Parse multipart form
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		logger.WarnContext(ctx, "Handler: Failed to parse multipart form",
			"error", err,
			"operation", "upload_piece")
		httpx.RespondWithError(w, fmt.Errorf("file too large or invalid form data"), http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.WarnContext(ctx, "Handler: No file provided",
			"error", err,
			"operation", "upload_piece")
		httpx.RespondWithError(w, fmt.Errorf("file is required"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	notes := r.FormValue("notes")

	logger.InfoContext(ctx, "Handler: Processing upload piece request",
		"operation", "upload_piece",
		"case_id", caseID.String(),
		"file_name", fileHeader.Filename,
		"file_size", fileHeader.Size)

	response, err := h.svc.UploadPiece(ctx, caseID, file, fileHeader, notes)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "upload piece")
		return
	}

	logger.InfoContext(ctx, "Handler: Upload piece completed successfully",
		"operation", "upload_piece",
		"piece_id", response.ID)

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
