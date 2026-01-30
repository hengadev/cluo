package ai

import "errors"

var (
	// ErrTransformationTypeInvalid is returned when an invalid transformation type is provided.
	ErrTransformationTypeInvalid = errors.New("invalid transformation type")

	// ErrTextTooLong is returned when input text exceeds maximum length.
	ErrTextTooLong = errors.New("text exceeds maximum length")

	// ErrEmptyInput is returned when input text is empty.
	ErrEmptyInput = errors.New("input cannot be empty")

	// ErrTranscriptionFailed is returned when transcription fails.
	ErrTranscriptionFailed = errors.New("transcription failed")

	// ErrInvalidAudioFormat is returned when audio format is invalid.
	ErrInvalidAudioFormat = errors.New("invalid audio format")

	// ErrAnalysisFailed is returned when analysis fails.
	ErrAnalysisFailed = errors.New("analysis failed")

	// ErrInvalidResponse is returned when AI response is invalid.
	ErrInvalidResponse = errors.New("invalid AI response")

	// ErrAIJobNotCancellable is returned when a job cannot be cancelled.
	ErrAIJobNotCancellable = errors.New("job cannot be cancelled (already processing or completed)")
)
