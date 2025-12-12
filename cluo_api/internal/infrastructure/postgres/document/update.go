package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// Update updates an existing document.
func (r *Repository) Update(ctx context.Context, doc document.Documentable) error {
	switch d := doc.(type) {
	case *document.EstimateEncx:
		return r.UpdateEstimate(ctx, d)
	case *document.MandateEncx:
		return r.UpdateMandate(ctx, d)
	case *document.ContractEncx:
		return r.UpdateContract(ctx, d)
	case *document.InvoiceEncx:
		return r.UpdateInvoice(ctx, d)
	default:
		return fmt.Errorf("unsupported document type: %T", doc)
	}
}
