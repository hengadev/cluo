package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/doctemplate"
	"github.com/hengadev/cluo_api/internal/common/pdf"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// RenderDocumentPDF fetches a document by ID, decrypts it, and renders it as a PDF.
func (s *Service) RenderDocumentPDF(ctx context.Context, id string, docType document.DocumentType) ([]byte, error) {
	encDoc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}

	decrypted, err := document.DecryptDocumentable(ctx, s.crypto, encDoc)
	if err != nil {
		return nil, fmt.Errorf("decrypt document: %w", err)
	}

	html, err := doctemplate.RenderDocument(decrypted)
	if err != nil {
		return nil, fmt.Errorf("render document html: %w", err)
	}

	pdfBytes, err := pdf.GenerateFromHTML(html)
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}

	return pdfBytes, nil
}
