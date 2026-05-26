package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Ensure SMTPAdapter implements EmailService at compile time.
var _ ports.EmailService = (*SMTPAdapter)(nil)

// SMTPAdapter sends emails via a standard SMTP server.
type SMTPAdapter struct {
	host     string
	port     string
	from     string
	username string
	password string
	log      *slog.Logger
}

// NewSMTPAdapter creates a new SMTP adapter from configuration.
func NewSMTPAdapter(cfg config.SMTPConfig, logger *slog.Logger) *SMTPAdapter {
	return &SMTPAdapter{
		host:     cfg.Host,
		port:     cfg.Port,
		from:     cfg.From,
		username: cfg.Username,
		password: cfg.Password,
		log:      logger,
	}
}

// Send delivers an HTML email to a single recipient.
func (a *SMTPAdapter) Send(ctx context.Context, to, subject, bodyHTML string) error {
	addr := a.host + ":" + a.port

	from := a.from
	recipients := []string{to}

	// Build the MIME message
	var msg strings.Builder
	msg.WriteString("From: " + from + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(bodyHTML)

	var auth smtp.Auth
	if a.username != "" && a.password != "" {
		auth = smtp.PlainAuth("", a.username, a.password, a.host)
	}

	if err := smtp.SendMail(addr, auth, from, recipients, []byte(msg.String())); err != nil {
		a.log.ErrorContext(ctx, "SMTP send failed",
			"error", err,
			"to", to,
			"subject", subject,
		)
		return fmt.Errorf("smtp send: %w", err)
	}

	a.log.InfoContext(ctx, "Email sent",
		"to", to,
		"subject", subject,
	)
	return nil
}
