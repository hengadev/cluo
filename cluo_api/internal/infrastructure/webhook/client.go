package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultTimeout is the default timeout for webhook requests.
	DefaultTimeout = 30 * time.Second
	// MaxRetries is the maximum number of retry attempts for webhooks.
	MaxRetries = 3
	// BaseRetryDelay is the base delay for exponential backoff.
	BaseRetryDelay = 1 * time.Second
)

// Client handles webhook notifications.
type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
}

// New creates a new webhook client.
func New(logger *slog.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		logger: logger,
	}
}

// WebhookPayload represents the payload sent to webhook URLs.
type WebhookPayload struct {
	Event           string    `json:"event"`
	JobID           uuid.UUID `json:"jobId"`
	Status          string    `json:"status"`
	TranscriptionID *uuid.UUID `json:"transcriptionId,omitempty"`
	ErrorMessage    string    `json:"errorMessage,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

// NotifyCompletion sends a completion notification to the webhook URL.
func (c *Client) NotifyCompletion(ctx context.Context, url string, jobID, transcriptionID uuid.UUID) error {
	payload := WebhookPayload{
		Event:           "transcription.completed",
		JobID:           jobID,
		Status:          "completed",
		TranscriptionID: &transcriptionID,
		Timestamp:       time.Now(),
	}

	return c.sendWithRetry(ctx, url, payload)
}

// NotifyFailure sends a failure notification to the webhook URL.
func (c *Client) NotifyFailure(ctx context.Context, url string, jobID uuid.UUID, errMsg string) error {
	payload := WebhookPayload{
		Event:        "transcription.failed",
		JobID:        jobID,
		Status:       "failed",
		ErrorMessage: errMsg,
		Timestamp:    time.Now(),
	}

	return c.sendWithRetry(ctx, url, payload)
}

// NotifyProgress sends a progress update to the webhook URL.
func (c *Client) NotifyProgress(ctx context.Context, url string, jobID uuid.UUID, progress int) error {
	payload := WebhookPayload{
		Event:     "transcription.progress",
		JobID:     jobID,
		Status:    "processing",
		Timestamp: time.Now(),
	}

	// We include progress in the event name for simplicity
	// Alternatively, we could add a Progress field to WebhookPayload
	c.logger.Debug("webhook progress notification",
		"url", url,
		"jobID", jobID,
		"progress", progress)

	return c.sendWithRetry(ctx, url, payload)
}

// sendWithRetry sends a webhook payload with retry logic.
func (c *Client) sendWithRetry(ctx context.Context, url string, payload WebhookPayload) error {
	var lastErr error

	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2^attempt * baseDelay (1s, 2s, 4s, ...)
			backoff := time.Duration(1<<uint(attempt)) * BaseRetryDelay
			c.logger.Debug("webhook retry",
				"url", url,
				"attempt", attempt+1,
				"backoff", backoff)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := c.send(ctx, url, payload); err != nil {
			lastErr = err
			c.logger.Warn("webhook attempt failed",
				"url", url,
				"attempt", attempt+1,
				"error", err)

			// Don't retry on client errors (4xx)
			if isClientError(err) {
				return err
			}

			continue
		}

		// Success
		return nil
	}

	// All retries exhausted
	return fmt.Errorf("webhook failed after %d attempts: %w", MaxRetries, lastErr)
}

// send sends a webhook payload.
func (c *Client) send(ctx context.Context, url string, payload WebhookPayload) error {
	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CLUO-Webhook/1.0")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &WebhookError{
			StatusCode: resp.StatusCode,
			URL:        url,
			Message:    fmt.Sprintf("webhook returned status %d", resp.StatusCode),
		}
	}

	c.logger.Debug("webhook sent successfully",
		"url", url,
		"event", payload.Event,
		"status", resp.StatusCode)

	return nil
}

// WebhookError represents an error from a webhook request.
type WebhookError struct {
	StatusCode int
	URL        string
	Message    string
}

func (e *WebhookError) Error() string {
	return fmt.Sprintf("webhook error: %s (status %d)", e.Message, e.StatusCode)
}

// isClientError returns true if the error is a client error (4xx).
func isClientError(err error) bool {
	if webhookErr, ok := err.(*WebhookError); ok {
		return webhookErr.StatusCode >= 400 && webhookErr.StatusCode < 500
	}
	return false
}
