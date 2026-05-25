package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/tiptap"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) GetReportHTMLByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")

	logger.InfoContext(ctx, "Handler: Processing get report HTML by token request",
		"operation", "get_report_html_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	// Validate token and retrieve caseID
	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "validate token for report HTML")
		return
	}

	// Retrieve the rapport (decrypted TipTap JSON)
	response, err := h.rapportSvc.GetRapportByCaseID(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get report for HTML conversion")
		return
	}

	// Convert TipTap JSON → HTML
	html, err := tiptap.ToHTML(response.Content)
	if err != nil {
		logger.ErrorContext(ctx, "Handler: Failed to convert TipTap JSON to HTML",
			"error", err.Error(),
			"operation", "get_report_html_by_token",
		)
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Get report HTML by token completed successfully",
		"operation", "get_report_html_by_token",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}
