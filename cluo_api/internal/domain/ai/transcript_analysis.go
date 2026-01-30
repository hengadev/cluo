package ai

import (
	"time"

	"github.com/google/uuid"
)

// TranscriptAnalysis represents an analysis of a transcription.
type TranscriptAnalysis struct {
	ID                  uuid.UUID              `db:"id"`
	TranscriptionID     uuid.UUID              `db:"transcription_id"`
	KeyFindings         string                 `encx:"encrypt" db:"key_findings_encrypted"`
	Summary             string                 `encx:"encrypt" db:"summary_encrypted"`
	Sentiment           string                 `db:"sentiment"`
	Topics              string                 `encx:"encrypt" db:"topics_encrypted"` // JSON array of topics
	SuggestedActions    string                 `encx:"encrypt" db:"suggested_actions_encrypted"`
	ModelUsed           string                 `db:"model_used"`
	ProcessingTimeMs    int64                  `db:"processing_time_ms"`
	CreatedAt           time.Time              `db:"created_at"`
}

// Sentiment represents the sentiment of analyzed text.
type Sentiment string

const (
	SentimentPositive   Sentiment = "positive"
	SentimentNeutral    Sentiment = "neutral"
	SentimentNegative   Sentiment = "negative"
	SentimentMixed      Sentiment = "mixed"
)

// Valid returns true if the sentiment is valid.
func (s Sentiment) Valid() bool {
	switch s {
	case SentimentPositive, SentimentNeutral, SentimentNegative, SentimentMixed:
		return true
	default:
		return false
	}
}

// NewTranscriptAnalysis creates a new TranscriptAnalysis entity.
func NewTranscriptAnalysis(
	transcriptionID uuid.UUID,
	keyFindings, summary string,
	sentiment string,
	topics, suggestedActions string,
	modelUsed string,
	processingTimeMs int64,
) *TranscriptAnalysis {
	return &TranscriptAnalysis{
		ID:               uuid.New(),
		TranscriptionID:  transcriptionID,
		KeyFindings:      keyFindings,
		Summary:          summary,
		Sentiment:        sentiment,
		Topics:           topics,
		SuggestedActions: suggestedActions,
		ModelUsed:        modelUsed,
		ProcessingTimeMs: processingTimeMs,
		CreatedAt:        time.Now(),
	}
}

// ReportSuggestion represents a suggested report based on transcript analysis.
type ReportSuggestion struct {
	ID              uuid.UUID `db:"id"`
	AnalysisID      uuid.UUID `db:"analysis_id"`
	ReportType      string    `db:"report_type"`
	Reasoning       string    `encx:"encrypt" db:"reasoning_encrypted"`
	Confidence      float32   `db:"confidence"`
	CreatedAt       time.Time `db:"created_at"`
}

// NewReportSuggestion creates a new ReportSuggestion entity.
func NewReportSuggestion(
	analysisID uuid.UUID,
	reportType, reasoning string,
	confidence float32,
) *ReportSuggestion {
	return &ReportSuggestion{
		ID:         uuid.New(),
		AnalysisID: analysisID,
		ReportType: reportType,
		Reasoning:  reasoning,
		Confidence: confidence,
		CreatedAt:  time.Now(),
	}
}
