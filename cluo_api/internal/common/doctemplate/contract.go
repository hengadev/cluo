package doctemplate

import (
	"fmt"
	"strings"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// RenderContract renders a Contract as HTML for PDF generation.
func RenderContract(c *document.Contract) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<h1>Contrat %s</h1>", pdfEscape(c.ContractNumber)))

	sb.WriteString(fmt.Sprintf("<p>Date de début : %s</p>", c.StartDate.Format("02/01/2006")))
	if c.EndDate != nil {
		sb.WriteString(fmt.Sprintf("<p>Date de fin : %s</p>", c.EndDate.Format("02/01/2006")))
	}

	sb.WriteString("<h2>Périmètre des prestations</h2>")
	sb.WriteString(fmt.Sprintf("<p>%s</p>", pdfEscape(c.ScopeOfServices)))

	sb.WriteString("<h2>Conditions de paiement</h2>")
	sb.WriteString(fmt.Sprintf("<p>%s</p>", pdfEscape(c.PaymentTerms)))

	if c.ContractValue != nil {
		sb.WriteString(fmt.Sprintf("<p>Valeur du contrat : %.2f", *c.ContractValue))
		if c.Currency != nil {
			sb.WriteString(fmt.Sprintf(" %s", pdfEscape(*c.Currency)))
		}
		sb.WriteString("</p>")
	}

	sb.WriteString("<h2>Clause de confidentialité</h2>")
	sb.WriteString(fmt.Sprintf("<p>%s</p>", pdfEscape(c.Confidentiality)))

	sb.WriteString("<h2>Clause de résiliation</h2>")
	sb.WriteString(fmt.Sprintf("<p>%s</p>", pdfEscape(c.TerminationClause)))

	if len(c.Signatures) > 0 {
		sb.WriteString("<h2>Signatures</h2>")
		for _, sig := range c.Signatures {
			sb.WriteString(fmt.Sprintf("<p>%s (%s) — signé le %s</p>",
				pdfEscape(sig.Name),
				pdfEscape(sig.Role),
				sig.SignedAt.Format("02/01/2006"),
			))
		}
	}

	return sb.String()
}
