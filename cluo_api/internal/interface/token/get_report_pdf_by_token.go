package tokenHandler

import (
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/common/pdf"
	"github.com/hengadev/cluo_api/internal/common/tiptap"
)

func (h *TokenHandler) GetReportPDFByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")

	logger.InfoContext(ctx, "Handler: Processing get report PDF by token request",
		"operation", "get_report_pdf_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	// Validate token and retrieve caseID
	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate token for report PDF")
		return
	}

	// Retrieve the rapport (decrypted TipTap JSON)
	response, err := h.rapportSvc.GetRapportByCaseID(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get report for PDF conversion")
		return
	}

	// Convert TipTap JSON → HTML
	html, err := tiptap.ToHTML(response.Content)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to convert TipTap JSON to HTML for PDF",
			"error", err.Error(),
			"operation", "get_report_pdf_by_token",
		)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Generate PDF from HTML
	pdfBytes, err := pdf.GenerateFromHTML(html)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to generate PDF from HTML",
			"error", err.Error(),
			"operation", "get_report_pdf_by_token",
		)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Get report PDF by token completed successfully",
		"operation", "get_report_pdf_by_token",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK,
		"pdf_size_bytes", len(pdfBytes))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="rapport.pdf"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	_, _ = w.Write(pdfBytes)
}
