package aiSpeechToTextHandler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
)

const (
	// maxMemSize is the maximum memory size for parsing multipart forms (100 MB).
	maxMemSize = 100 << 20 // 100 MB
)

type SubmitJobResponse struct {
	JobID      uuid.UUID  `json:"jobId"`
	MediaFileID uuid.UUID `json:"mediaFileId"`
	Status     ai.JobStatus `json:"status"`
	Progress   int        `json:"progress"`
	CreatedAt  string     `json:"createdAt"`
}

func (h *handler) SubmitJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse multipart form (max 100 MB)
	if err := r.ParseMultipartForm(maxMemSize); err != nil {
		logger.WarnContext(ctx, "Handler: Failed to parse multipart form",
			"error", err,
			"operation", "submit_transcription_job",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.WarnContext(ctx, "Handler: Missing file in form",
			"error", err,
			"operation", "submit_transcription_job")
		httpx.RespondWithError(w, fmt.Errorf("file field is required"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	audioData, err := io.ReadAll(file)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Failed to read file",
			"error", err,
			"operation", "submit_transcription_job")
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Get media file ID from form
	mediaFileIDStr := r.FormValue("mediaFileId")
	if mediaFileIDStr == "" {
		httpx.RespondWithError(w, fmt.Errorf("mediaFileId is required"), http.StatusBadRequest)
		return
	}
	mediaFileID, err := uuid.Parse(mediaFileIDStr)
	if err != nil {
		httpx.RespondWithError(w, fmt.Errorf("invalid mediaFileId: %w", err), http.StatusBadRequest)
		return
	}

	// Get optional webhook URL
	var webhookURL *string
	if webhookURLStr := r.FormValue("webhookUrl"); webhookURLStr != "" {
		webhookURL = &webhookURLStr
	}

	// Get user ID from session context
	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		logger.ErrorContext(ctx, "Handler: Session info not found in context",
			"operation", "submit_transcription_job")
		httpx.RespondWithError(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}
	userID := sessionInfo.UserID

	logger.InfoContext(ctx, "Handler: Processing submit transcription job request",
		"operation", "submit_transcription_job",
		"method", r.Method,
		"path", r.URL.Path,
		"mediaFileID", mediaFileID,
		"filename", header.Filename,
		"fileSize", len(audioData),
		"userID", userID)

	// Call service
	req := &ports.SubmitTranscriptionRequest{
		MediaFileID:    mediaFileID,
		AudioData:      audioData,
		AudioFilename:  header.Filename,
		WebhookURL:     webhookURL,
		UserID:         userID,
	}

	job, err := h.svc.SubmitTranscriptionJob(ctx, req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "submit transcription job")
		return
	}

	response := &SubmitJobResponse{
		JobID:      job.ID,
		MediaFileID: job.MediaFileID,
		Status:     job.Status,
		Progress:   job.Progress,
		CreatedAt:  job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	logger.InfoContext(ctx, "Handler: Submit transcription job request completed successfully",
		"operation", "submit_transcription_job",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusAccepted,
		"jobID", job.ID)

	httpx.RespondWithJSON(w, response, http.StatusAccepted)
}
