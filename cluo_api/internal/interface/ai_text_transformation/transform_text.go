package aiTextTransformationHandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
)

type TransformTextRequest struct {
	InputText          string                 `json:"inputText"`
	TransformationType ai.TextTransformationType `json:"transformationType"`
}

type TransformTextResponse struct {
	ID                 uuid.UUID              `json:"id"`
	InputText          string                 `json:"inputText"`
	OutputText         string                 `json:"outputText"`
	TransformationType  ai.TextTransformationType `json:"transformationType"`
	ModelUsed          string                 `json:"modelUsed"`
	InputLength        int                    `json:"inputLength"`
	OutputLength       int                    `json:"outputLength"`
	ProcessingTimeMs   int64                  `json:"processingTimeMs"`
	CreatedAt          string                 `json:"createdAt"`
}

func (h *handler) TransformText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	var payload TransformTextRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "transform_text",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing transform text request",
		"operation", "transform_text",
		"method", r.Method,
		"path", r.URL.Path,
		"transformationType", payload.TransformationType,
		"inputLength", len(payload.InputText))

	// Call service
	req := &ports.TransformTextRequest{
		InputText:          payload.InputText,
		TransformationType: payload.TransformationType,
	}

	result, err := h.svc.TransformText(ctx, req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "transform text")
		return
	}

	response := &TransformTextResponse{
		ID:                result.ID,
		InputText:         result.InputText,
		OutputText:        result.OutputText,
		TransformationType: result.TransformationType,
		ModelUsed:         result.ModelUsed,
		InputLength:       result.InputLength,
		OutputLength:      result.OutputLength,
		ProcessingTimeMs:  result.ProcessingTimeMs,
		CreatedAt:         result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	logger.InfoContext(ctx, "Handler: Transform text request completed successfully",
		"operation", "transform_text",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
