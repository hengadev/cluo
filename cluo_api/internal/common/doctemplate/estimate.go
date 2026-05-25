package doctemplate

import (
	"fmt"
	"strings"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// RenderEstimate renders an Estimate as HTML for PDF generation.
func RenderEstimate(e *document.Estimate) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<h1>Devis %s</h1>", pdfEscape(e.EstimateNumber)))

	sb.WriteString(fmt.Sprintf("<p>Date d'émission : %s</p>", e.IssueDate.Format("02/01/2006")))

	if e.ValidUntil != nil {
		sb.WriteString(fmt.Sprintf("<p>Valide jusqu'au : %s</p>", e.ValidUntil.Format("02/01/2006")))
	}

	sb.WriteString("<h2>Lignes</h2>")
	for _, item := range e.LineItems {
		sb.WriteString(fmt.Sprintf("<p>%s — %.2f x %.2f = %.2f</p>",
			pdfEscape(item.Description),
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
		))
	}

	sb.WriteString(fmt.Sprintf("<h2>Total estimé : %.2f</h2>", e.EstimatedTotal))

	if e.Notes != nil && *e.Notes != "" {
		sb.WriteString(fmt.Sprintf("<p>Notes : %s</p>", pdfEscape(*e.Notes)))
	}

	return sb.String()
}

// pdfEscape escapes special characters for HTML output.
func pdfEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
