package tokenHandler

import (
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/doctemplate"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/common/pdf"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

var validDocTypes = map[string]document.DocumentType{
	"estimate": document.DocumentTypeEstimate,
	"mandate":  document.DocumentTypeMandate,
	"contract": document.DocumentTypeContract,
	"invoice":  document.DocumentTypeInvoice,
}

func (h *handler) GetDocumentPDFByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")
	docTypeStr := r.PathValue("type")

	logger.InfoContext(ctx, "Handler: Processing get document PDF by token request",
		"operation", "get_document_pdf_by_token",
		"method", r.Method,
		"path", r.URL.Path,
		"doc_type", docTypeStr)

	// Validate document type
	docType, ok := validDocTypes[docTypeStr]
	if !ok {
		httpx.RespondWithError(w, fmt.Errorf("unknown document type: %s", docTypeStr), http.StatusBadRequest)
		return
	}

	// Validate token and retrieve caseID
	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate token for document PDF")
		return
	}

	// Document repo is nil while document packages are disabled.
	if h.documentRepo == nil || h.crypto == nil {
		httpx.RespondWithError(w, fmt.Errorf("document feature not available"), http.StatusNotFound)
		return
	}

	// Fetch the encrypted document for this case + type
	encDoc, err := h.documentRepo.GetFirstByCaseAndType(ctx, caseID.String(), docType)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get document for PDF")
		return
	}

	// Decrypt the document
	decrypted, err := document.DecryptDocumentable(ctx, h.crypto, encDoc)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to decrypt document for PDF",
			"error", err.Error(),
			"operation", "get_document_pdf_by_token",
		)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Render HTML template
	html, err := doctemplate.RenderDocument(decrypted)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to render document HTML",
			"error", err.Error(),
			"operation", "get_document_pdf_by_token",
		)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Generate PDF from HTML
	pdfBytes, err := pdf.GenerateFromHTML(html)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to generate PDF from HTML",
			"error", err.Error(),
			"operation", "get_document_pdf_by_token",
		)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Get document PDF by token completed successfully",
		"operation", "get_document_pdf_by_token",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"pdf_size_bytes", len(pdfBytes))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, docTypeStr))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	_, _ = w.Write(pdfBytes)
}
