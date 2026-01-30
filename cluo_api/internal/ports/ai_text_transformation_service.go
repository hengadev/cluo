package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// TextTransformationService defines the interface for text transformation operations.
type TextTransformationService interface {
	// TransformText applies a transformation to input text.
	TransformText(ctx context.Context, req *TransformTextRequest) (*ai.TextTransformation, error)
}

// TransformTextRequest defines the request for text transformation.
type TransformTextRequest struct {
	InputText           string
	TransformationType  ai.TextTransformationType
}
