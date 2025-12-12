package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// Create creates a new document of any type.
func (r *Repository) Create(ctx context.Context, doc document.Documentable) error {
	switch d := doc.(type) {
	case *document.EstimateEncx:
		return r.CreateEstimate(ctx, d)
	case *document.MandateEncx:
		return r.CreateMandate(ctx, d)
	case *document.ContractEncx:
		return r.CreateContract(ctx, d)
	case *document.InvoiceEncx:
		return r.CreateInvoice(ctx, d)
	default:
		return fmt.Errorf("unsupported document type: %T", doc)
	}
}
