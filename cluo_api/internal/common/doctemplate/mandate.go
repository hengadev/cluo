package doctemplate

import (
	"fmt"
	"strings"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// RenderMandate renders a Mandate as HTML for PDF generation.
func RenderMandate(m *document.Mandate) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<h1>Mandat %s</h1>", pdfEscape(m.MandateNumber)))

	sb.WriteString(fmt.Sprintf("<p>Date d'émission : %s</p>", m.IssueDate.Format("02/01/2006")))

	sb.WriteString("<h2>Périmètre de la mission</h2>")
	sb.WriteString(fmt.Sprintf("<p>%s</p>", pdfEscape(m.ScopeOfWork)))

	sb.WriteString(fmt.Sprintf("<p>Validité : du %s", m.ValidFrom.Format("02/01/2006")))
	if m.ValidUntil != nil {
		sb.WriteString(fmt.Sprintf(" au %s", m.ValidUntil.Format("02/01/2006")))
	}
	sb.WriteString("</p>")

	sb.WriteString("<h2>Conditions</h2>")
	sb.WriteString(fmt.Sprintf("<p>%s</p>", pdfEscape(m.TermsConditions)))

	if m.SpecialInstructions != nil && *m.SpecialInstructions != "" {
		sb.WriteString(fmt.Sprintf("<p>Instructions spéciales : %s</p>", pdfEscape(*m.SpecialInstructions)))
	}

	sb.WriteString("<h2>Signatures</h2>")
	if m.ClientSignature != nil {
		sb.WriteString(fmt.Sprintf("<p>Client : %s (%s) — signé le %s</p>",
			pdfEscape(m.ClientSignature.Name),
			pdfEscape(m.ClientSignature.Role),
			m.ClientSignature.SignedAt.Format("02/01/2006"),
		))
	}
	if m.InvestigatorSignature != nil {
		sb.WriteString(fmt.Sprintf("<p>Enquêteur : %s (%s) — signé le %s</p>",
			pdfEscape(m.InvestigatorSignature.Name),
			pdfEscape(m.InvestigatorSignature.Role),
			m.InvestigatorSignature.SignedAt.Format("02/01/2006"),
		))
	}

	return sb.String()
}
