package document

import (
	"context"
	"fmt"

	"github.com/hengadev/encx"
)

// EncryptDocumentable dispatches to the appropriate Process<Type>Encx function
// based on the concrete type of doc.
func EncryptDocumentable(ctx context.Context, crypto encx.CryptoService, doc interface{}) (Documentable, error) {
	switch d := doc.(type) {
	case *Estimate:
		return ProcessEstimateEncx(ctx, crypto, d)
	case *Mandate:
		return ProcessMandateEncx(ctx, crypto, d)
	case *Contract:
		return ProcessContractEncx(ctx, crypto, d)
	case *Invoice:
		return ProcessInvoiceEncx(ctx, crypto, d)
	default:
		return nil, fmt.Errorf("document: cannot encrypt unknown type %T", doc)
	}
}
