package ai

import (
	"time"

	"github.com/google/uuid"
)

// TextTransformation represents a text transformation result.
type TextTransformation struct {
	ID                  uuid.UUID              `db:"id"`
	InputText           string                 `encx:"encrypt" db:"input_text_encrypted"`
	OutputText          string                 `encx:"encrypt" db:"output_text_encrypted"`
	TransformationType  TextTransformationType `db:"transformation_type"`
	ModelUsed           string                 `db:"model_used"`
	InputLength         int                    `db:"input_length"`
	OutputLength        int                    `db:"output_length"`
	ProcessingTimeMs    int64                  `db:"processing_time_ms"`
	CreatedAt           time.Time              `db:"created_at"`
}

// NewTextTransformation creates a new TextTransformation entity.
func NewTextTransformation(
	inputText, outputText string,
	transformationType TextTransformationType,
	modelUsed string,
	processingTimeMs int64,
) *TextTransformation {
	return &TextTransformation{
		ID:                  uuid.New(),
		InputText:           inputText,
		OutputText:          outputText,
		TransformationType:  transformationType,
		ModelUsed:           modelUsed,
		InputLength:         len(inputText),
		OutputLength:        len(outputText),
		ProcessingTimeMs:    processingTimeMs,
		CreatedAt:           time.Now(),
	}
}
