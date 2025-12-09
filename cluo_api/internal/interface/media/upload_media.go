package mediaHandler

import (
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

const (
	maxUploadSize = 100 << 20 // 100 MB
)

func (h *handler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse multipart form
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		logger.WarnContext(ctx, "Handler: Failed to parse multipart form",
			"error", err,
			"operation", "upload_media",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, fmt.Errorf("file too large or invalid form data"), http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.WarnContext(ctx, "Handler: No file provided",
			"error", err,
			"operation", "upload_media")
		httpx.RespondWithError(w, fmt.Errorf("file is required"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get metadata from form
	caseID := r.FormValue("caseId")
	caption := r.FormValue("caption")
	isPublishedStr := r.FormValue("isPublished")

	// Get MIME type
	mimeType := fileHeader.Header.Get("Content-Type")

	// Build UploadMediaRequest
	var captionPtr *string
	if caption != "" {
		captionPtr = &caption
	}

	var isPublishedPtr *bool
	if isPublishedStr != "" {
		isPublished := isPublishedStr == "true"
		isPublishedPtr = &isPublished
	}

	request := &domain.UploadMediaRequest{
		CaseID:      caseID,
		File:        file,
		FileName:    fileHeader.Filename,
		MimeType:    mimeType,
		FileSize:    fileHeader.Size,
		Caption:     captionPtr,
		IsPublished: isPublishedPtr,
	}

	logger.InfoContext(ctx, "Handler: Processing upload media request",
		"operation", "upload_media",
		"case_id", caseID,
		"file_name", fileHeader.Filename,
		"file_size", fileHeader.Size,
		"mime_type", mimeType)

	response, err := h.svc.UploadMedia(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "upload media")
		return
	}

	logger.InfoContext(ctx, "Handler: Upload media completed successfully",
		"operation", "upload_media",
		"media_id", response.ID)

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}
