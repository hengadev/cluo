package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// SendDocument sends a document to recipients.
func (s *Service) SendDocument(ctx context.Context, id string, docType document.DocumentType, req *document.SendDocumentRequest) error {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	// State machine enforcement: only draft documents can be sent
	if err := s.validateDocumentTransition(doc, document.DocumentStatusSent); err != nil {
		return err
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
			"error", err,
			"document_id", id,
		)
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

	// Send to all recipients in background
	for _, to := range req.Recipients {
		go func(recipient string) {
			bgCtx := context.Background()
			if err := s.emailService.Send(bgCtx, recipient, subject, bodyHTML); err != nil {
				s.logger.ErrorContext(bgCtx, "Failed to send document email",
					"error", err,
					"to", recipient,
					"document_id", id,
				)
			}
		}(to)
	}

	return nil
}