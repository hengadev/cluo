package doctemplate

import (
	"fmt"
	"strings"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// RenderInvoice renders an Invoice as HTML for PDF generation.
func RenderInvoice(i *document.Invoice) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<h1>Facture %s</h1>", pdfEscape(i.InvoiceNumber)))

	sb.WriteString(fmt.Sprintf("<p>Date d'émission : %s</p>", i.IssueDate.Format("02/01/2006")))
	sb.WriteString(fmt.Sprintf("<p>Date d'échéance : %s</p>", i.DueDate.Format("02/01/2006")))

	sb.WriteString("<h2>Lignes</h2>")
	for _, item := range i.LineItems {
		sb.WriteString(fmt.Sprintf("<p>%s — %.2f x %.2f = %.2f</p>",
			pdfEscape(item.Description),
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
		))
	}

	subtotal := i.TotalAmount - i.TaxAmount
	sb.WriteString(fmt.Sprintf("<p>Sous-total : %.2f</p>", subtotal))
	sb.WriteString(fmt.Sprintf("<p>Tax (%.0f%%) : %.2f</p>", i.TaxRate*100, i.TaxAmount))
	sb.WriteString(fmt.Sprintf("<h2>Total : %.2f</h2>", i.TotalAmount))

	sb.WriteString(fmt.Sprintf("<p>Statut de paiement : %s</p>", string(i.PaymentStatus)))

	if i.Notes != nil && *i.Notes != "" {
		sb.WriteString(fmt.Sprintf("<p>Notes : %s</p>", pdfEscape(*i.Notes)))
	}

	return sb.String()
}
