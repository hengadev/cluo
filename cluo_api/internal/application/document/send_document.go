package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// SendDocument sends a document to recipients.
// If req.Recipients is empty and req.SendEmail is true, the service automatically
// resolves the recipient list from the document's client primary contacts.
func (s *Service) SendDocument(ctx context.Context, id string, docType document.DocumentType, req *document.SendDocumentRequest) error {
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	// State machine enforcement: only draft documents can be sent
	if err := s.validateDocumentTransition(doc, document.DocumentStatusSent); err != nil {
		return err
	}

	// Resolve recipients: if none supplied, fall back to the client's contacts.
	// All concrete document types embed DocumentBase which exposes GetClientID();
	// the *Encx types do not, so we guard with an interface assertion.
	type hasClientID interface {
		GetClientID() uuid.UUID
	}
	recipients := req.Recipients
	if len(recipients) == 0 && req.SendEmail {
		if withClient, ok := doc.(hasClientID); ok {
			resolved, resolveErr := s.resolveClientEmails(ctx, withClient.GetClientID())
			if resolveErr != nil {
				s.logger.WarnContext(ctx, "Could not resolve client emails for document send",
					"error", resolveErr, "document_id", id)
			}
			recipients = resolved
		}
	}

	// Update document status to Sent
	doc.SetStatus(document.DocumentStatusSent)
	doc.UpdateTimestamp()

	if err := s.repo.Update(ctx, doc); err != nil {
		return errs.NewNotUpdatedErr(err, "document")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, doc, &authorID, stringPtr("Document sent")); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create document version",
			"error", err, "document_id", id)
	}

	if len(recipients) == 0 {
		return nil
	}

	// Dispatch email notification asynchronously (fire-and-forget).
	// Build email body before spawning goroutine to avoid data races on doc.
	docTypeLabel := string(docType)
	var docNumber string
	switch d := doc.(type) {
	case *document.Estimate:
		docNumber = d.EstimateNumber
	case *document.Mandate:
		docNumber = d.MandateNumber
	case *document.Contract:
		docNumber = d.ContractNumber
	case *document.Invoice:
		docNumber = d.InvoiceNumber
	}

	subject := fmt.Sprintf("Document envoyé : %s %s", docTypeLabel, docNumber)
	bodyHTML := fmt.Sprintf(`
		<html><body>
		<p>Bonjour,</p>
		<p>Un document a été envoyé :</p>
		<ul>
			<li><strong>Type :</strong> %s</li>
			<li><strong>Numéro :</strong> %s</li>
		</ul>
		<p>Ce document a été marqué comme envoyé.</p>
		</body></html>
	`, docTypeLabel, docNumber)

	for _, to := range recipients {
		go func(recipient string) {
			bgCtx := context.Background()
			if err := s.emailService.Send(bgCtx, recipient, subject, bodyHTML); err != nil {
				s.logger.ErrorContext(bgCtx, "Failed to send document email",
					"error", err, "to", recipient, "document_id", id)
			}
		}(to)
	}

	return nil
}

// resolveClientEmails returns the decrypted email addresses for all contacts of a client.
func (s *Service) resolveClientEmails(ctx context.Context, clientID uuid.UUID) ([]string, error) {
	contactEncxs, err := s.clientRepo.GetAllContactsByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch contacts for client %s: %w", clientID, err)
	}

	var emails []string
	for _, encx := range contactEncxs {
		contact, err := client.DecryptContactEncx(ctx, s.crypto, encx)
		if err != nil {
			s.logger.WarnContext(ctx, "Failed to decrypt contact", "error", err)
			continue
		}
		if contact.Email != "" {
			emails = append(emails, contact.Email)
		}
	}
	return emails, nil
}