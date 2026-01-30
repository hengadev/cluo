package errs

import "errors"

var (
	// AI Model errors
	ErrAIModelUnavailable   = errors.New("ai model unavailable")
	ErrAITranscriptionFailed = errors.New("transcription failed")
	ErrAIAnalysisFailed      = errors.New("analysis failed")
	ErrAIInvalidAudioFormat  = errors.New("invalid audio format")
	ErrAITextTooLong         = errors.New("text exceeds maximum length")
	ErrAIEmptyInput          = errors.New("input cannot be empty")
	ErrAIInvalidResponse     = errors.New("invalid ai response")
	ErrAITimeout             = errors.New("ai request timeout")
	ErrAILocalhostRequired   = errors.New("ai service must be localhost only")

	// Job-related errors
	ErrAIJobNotFound       = errors.New("transcription job not found")
	ErrAIJobNotCancellable = errors.New("job cannot be cancelled (already processing or completed)")
	ErrAIJobAlreadyClaimed = errors.New("job already claimed by another worker")
	ErrAIQueueFull         = errors.New("transcription queue is full, try again later")
)

// NewAIModelUnavailableErr creates an error for unavailable AI model.
func NewAIModelUnavailableErr(message string) error {
	return wrapAIError(ErrAIModelUnavailable, message)
}

// NewAITranscriptionFailedErr creates an error for failed transcription.
func NewAITranscriptionFailedErr(err error) error {
	return wrapAIError(ErrAITranscriptionFailed, err.Error())
}

// NewAIAnalysisFailedErr creates an error for failed analysis.
func NewAIAnalysisFailedErr(err error) error {
	return wrapAIError(ErrAIAnalysisFailed, err.Error())
}

// NewAIInvalidAudioFormatErr creates an error for invalid audio format.
func NewAIInvalidAudioFormatErr(format string) error {
	return wrapAIError(ErrAIInvalidAudioFormat, format)
}

// NewAITextTooLongErr creates an error for text that exceeds maximum length.
func NewAITextTooLongErr(maxLen int) error {
	return wrapAIError(ErrAITextTooLong, "maximum length is "+string(rune(maxLen)))
}

// NewAIEmptyInputErr creates an error for empty input.
func NewAIEmptyInputErr() error {
	return ErrAIEmptyInput
}

// NewAIInvalidResponseErr creates an error for invalid AI response.
func NewAIInvalidResponseErr(reason string) error {
	return wrapAIError(ErrAIInvalidResponse, reason)
}

// NewAITimeoutErr creates an error for AI request timeout.
func NewAITimeoutErr() error {
	return ErrAITimeout
}

// NewAILocalhostRequiredErr creates an error for non-localhost AI service.
func NewAILocalhostRequiredErr() error {
	return ErrAILocalhostRequired
}

// NewAIJobNotFoundErr creates an error for job not found.
func NewAIJobNotFoundErr() error {
	return ErrAIJobNotFound
}

// NewAIJobNotCancellableErr creates an error for non-cancellable job.
func NewAIJobNotCancellableErr(status string) error {
	return wrapAIError(ErrAIJobNotCancellable, "current status is "+status)
}

// NewAIJobAlreadyClaimedErr creates an error for already claimed job.
func NewAIJobAlreadyClaimedErr() error {
	return ErrAIJobAlreadyClaimed
}

// NewAIQueueFullErr creates an error for full queue.
func NewAIQueueFullErr() error {
	return ErrAIQueueFull
}

func wrapAIError(base error, message string) error {
	if message == "" {
		return base
	}
	return &wrappedAIError{
		base:    base,
		message: message,
	}
}

type wrappedAIError struct {
	base    error
	message string
}

func (e *wrappedAIError) Error() string {
	return e.base.Error() + ": " + e.message
}

func (e *wrappedAIError) Unwrap() error {
	return e.base
}
