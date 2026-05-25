package doctemplate

import (
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// RenderDocument renders a decrypted document struct into HTML suitable for
// the PDF generator. It dispatches by concrete type.
func RenderDocument(doc interface{}) (string, error) {
	switch d := doc.(type) {
	case *document.Estimate:
		return RenderEstimate(d), nil
	case *document.Mandate:
		return RenderMandate(d), nil
	case *document.Contract:
		return RenderContract(d), nil
	case *document.Invoice:
		return RenderInvoice(d), nil
	default:
		return "", fmt.Errorf("doctemplate: unsupported document type %T", doc)
	}
}
