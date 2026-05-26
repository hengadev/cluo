package ports

import "context"

// EmailService is a port for sending email notifications.
type EmailService interface {
	// Send dispatches an HTML email to a single recipient.
	// Implementations must be safe to call from goroutines.
	Send(ctx context.Context, to, subject, bodyHTML string) error
}
