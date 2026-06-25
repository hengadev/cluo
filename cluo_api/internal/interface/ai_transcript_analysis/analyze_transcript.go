package aiTranscriptAnalysisHandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
)

type AnalyzeTranscriptRequest struct {
	TranscriptionID string `json:"transcriptionId"`
}

type AnalyzeTranscriptResponse struct {
	ID               uuid.UUID `json:"id"`
	TranscriptionID  uuid.UUID `json:"transcriptionId"`
	KeyFindings      string    `json:"keyFindings"`
	Summary          string    `json:"summary"`
	Sentiment        string    `json:"sentiment"`
	Topics           string    `json:"topics"` // JSON array as string
	SuggestedActions string    `json:"suggestedActions"`
	ModelUsed        string    `json:"modelUsed"`
	ProcessingTimeMs int64     `json:"processingTimeMs"`
	CreatedAt        string    `json:"createdAt"`
}

func (h *handler) AnalyzeTranscript(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload AnalyzeTranscriptRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "analyze_transcript",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	transcriptionID, err := uuid.Parse(payload.TranscriptionID)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing analyze transcript request",
		"operation", "analyze_transcript",
		"method", r.Method,
		"path", r.URL.Path,
		"transcriptionID", transcriptionID)

	// Call service
	req := &ports.AnalyzeTranscriptRequest{
		TranscriptionID: transcriptionID,
	}

	result, err := h.svc.AnalyzeTranscript(ctx, req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "analyze transcript")
		return
	}

	response := &AnalyzeTranscriptResponse{
		ID:               result.ID,
		TranscriptionID:  result.TranscriptionID,
		KeyFindings:      result.KeyFindings,
		Summary:          result.Summary,
		Sentiment:        result.Sentiment,
		Topics:           result.Topics,
		SuggestedActions: result.SuggestedActions,
		ModelUsed:        result.ModelUsed,
		ProcessingTimeMs: result.ProcessingTimeMs,
		CreatedAt:        result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	logger.InfoContext(ctx, "Handler: Analyze transcript request completed successfully",
		"operation", "analyze_transcript",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"analysisID", result.ID)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

type GetAnalysisResponse struct {
	ID               uuid.UUID `json:"id"`
	TranscriptionID  uuid.UUID `json:"transcriptionId"`
	KeyFindings      string    `json:"keyFindings"`
	Summary          string    `json:"summary"`
	Sentiment        ai.Sentiment `json:"sentiment"`
	Topics           string    `json:"topics"`
	SuggestedActions string    `json:"suggestedActions"`
	CreatedAt        string    `json:"createdAt"`
}

func (h *handler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract analysis ID from URL path
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.DebugContext(ctx, "Handler: Processing get analysis request",
		"operation", "get_analysis",
		"analysisID", id)

	// Call service
	result, err := h.svc.GetAnalysis(ctx, id)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get analysis")
		return
	}

	response := &GetAnalysisResponse{
		ID:               result.ID,
		TranscriptionID:  result.TranscriptionID,
		KeyFindings:      result.KeyFindings,
		Summary:          result.Summary,
		Sentiment:        ai.Sentiment(result.Sentiment),
		Topics:           result.Topics,
		SuggestedActions: result.SuggestedActions,
		CreatedAt:        result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

// GetAnalysisByTranscriptionID retrieves an analysis by its transcription ID.
// Returns 404 when no analysis exists for the given transcription yet.
func (h *handler) GetAnalysisByTranscriptionID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract transcription ID from URL path
	transcriptionIDStr := r.PathValue("transcriptionId")
	transcriptionID, err := uuid.Parse(transcriptionIDStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.DebugContext(ctx, "Handler: Processing get analysis by transcription request",
		"operation", "get_analysis_by_transcription",
		"transcriptionID", transcriptionID)

	// Call service
	result, err := h.svc.GetAnalysisByTranscriptionID(ctx, transcriptionID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get analysis by transcription")
		return
	}

	response := &AnalyzeTranscriptResponse{
		ID:               result.ID,
		TranscriptionID:  result.TranscriptionID,
		KeyFindings:      result.KeyFindings,
		Summary:          result.Summary,
		Sentiment:        result.Sentiment,
		Topics:           result.Topics,
		SuggestedActions: result.SuggestedActions,
		ModelUsed:        result.ModelUsed,
		ProcessingTimeMs: result.ProcessingTimeMs,
		CreatedAt:        result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
