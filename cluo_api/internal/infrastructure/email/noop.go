package email

import (
	"context"
	"log/slog"

	"github.com/hengadev/cluo_api/internal/ports"
)

// Ensure NoOpAdapter implements EmailService at compile time.
var _ ports.EmailService = (*NoOpAdapter)(nil)

// NoOpAdapter discards all email sends. Used in development and tests.
type NoOpAdapter struct {
	log *slog.Logger
}

// NewNoOpAdapter creates a NoOp email service.
func NewNoOpAdapter(logger *slog.Logger) *NoOpAdapter {
	return &NoOpAdapter{log: logger}
}

// Send discards the email and returns nil.
func (n *NoOpAdapter) Send(ctx context.Context, to, subject, bodyHTML string) error {
	n.log.DebugContext(ctx, "NoOp email: discarding send",
		"to", to,
		"subject", subject,
	)
	return nil
}
