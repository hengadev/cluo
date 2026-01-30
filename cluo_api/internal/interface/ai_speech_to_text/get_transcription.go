package aiSpeechToTextHandler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

type TranscriptionResponse struct {
	ID               uuid.UUID `json:"id"`
	JobID            uuid.UUID `json:"jobId"`
	MediaFileID      uuid.UUID `json:"mediaFileId"`
	AudioURL         string    `json:"audioUrl"`
	Transcript       string    `json:"transcript"`
	ConfidenceScore  float32   `json:"confidenceScore"`
	Language         string    `json:"language"`
	Duration         int64     `json:"duration"` // milliseconds
	ModelName        string    `json:"modelName"`
	ProcessingTimeMs int64     `json:"processingTimeMs"`
	CreatedAt        string    `json:"createdAt"`
}

func (h *handler) GetTranscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract transcription ID from URL path
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.DebugContext(ctx, "Handler: Processing get transcription request",
		"operation", "get_transcription",
		"transcriptionID", id)

	// Call service
	transcription, err := h.svc.GetTranscription(ctx, id)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get transcription")
		return
	}

	response := &TranscriptionResponse{
		ID:               transcription.ID,
		JobID:            transcription.JobID,
		MediaFileID:      transcription.MediaFileID,
		AudioURL:         transcription.AudioURL,
		Transcript:       transcription.Transcript,
		ConfidenceScore:  transcription.ConfidenceScore,
		Language:         transcription.Language,
		Duration:         transcription.Duration,
		ModelName:        transcription.ModelName,
		ProcessingTimeMs: transcription.ProcessingTimeMs,
		CreatedAt:        transcription.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) DeleteTranscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract transcription ID from URL path
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing delete transcription request",
		"operation", "delete_transcription",
		"transcriptionID", id)

	// Call service
	if err := h.svc.DeleteTranscription(ctx, id); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete transcription")
		return
	}

	logger.InfoContext(ctx, "Handler: Delete transcription request completed successfully",
		"operation", "delete_transcription",
		"transcriptionID", id)

	httpx.RespondWithJSON(w, map[string]string{"message": "Transcription deleted"}, http.StatusOK)
}

type ListTranscriptionsResponse struct {
	Transcriptions []TranscriptionResponse `json:"transcriptions"`
	Total          int                     `json:"total"`
}

func (h *handler) ListTranscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract media file ID from query parameter
	mediaFileIDStr := r.URL.Query().Get("mediaFileId")
	if mediaFileIDStr == "" {
		httpx.RespondWithError(w, fmt.Errorf("mediaFileId is required"), http.StatusBadRequest)
		return
	}

	mediaFileID, err := uuid.Parse(mediaFileIDStr)
	if err != nil {
		httpx.RespondWithError(w, fmt.Errorf("invalid mediaFileId: %w", err), http.StatusBadRequest)
		return
	}

	logger.DebugContext(ctx, "Handler: Processing list transcriptions request",
		"operation", "list_transcriptions",
		"mediaFileID", mediaFileID)

	// Call service
	transcription, err := h.svc.GetTranscriptionByMediaFileID(ctx, mediaFileID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get transcription by media file")
		return
	}

	response := &TranscriptionResponse{
		ID:               transcription.ID,
		JobID:            transcription.JobID,
		MediaFileID:      transcription.MediaFileID,
		AudioURL:         transcription.AudioURL,
		Transcript:       transcription.Transcript,
		ConfidenceScore:  transcription.ConfidenceScore,
		Language:         transcription.Language,
		Duration:         transcription.Duration,
		ModelName:        transcription.ModelName,
		ProcessingTimeMs: transcription.ProcessingTimeMs,
		CreatedAt:        transcription.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	listResponse := ListTranscriptionsResponse{
		Transcriptions: []TranscriptionResponse{*response},
		Total:          1,
	}

	httpx.RespondWithJSON(w, listResponse, http.StatusOK)
}
