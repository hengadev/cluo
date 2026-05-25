package document

import (
	"context"
	"fmt"

	"github.com/hengadev/encx"
)

// DecryptDocumentable dispatches to the appropriate DecryptXxxEncx function
// based on the concrete type of encDoc.
func DecryptDocumentable(ctx context.Context, crypto encx.CryptoService, encDoc Documentable) (interface{}, error) {
	switch d := encDoc.(type) {
	case *EstimateEncx:
		return DecryptEstimateEncx(ctx, crypto, d)
	case *MandateEncx:
		return DecryptMandateEncx(ctx, crypto, d)
	case *ContractEncx:
		return DecryptContractEncx(ctx, crypto, d)
	case *InvoiceEncx:
		return DecryptInvoiceEncx(ctx, crypto, d)
	default:
		return nil, fmt.Errorf("document: cannot decrypt unknown type %T", encDoc)
	}
}
