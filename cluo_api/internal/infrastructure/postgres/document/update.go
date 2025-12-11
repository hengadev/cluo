package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// Update updates an existing document.
func (r *Repository) Update(ctx context.Context, doc document.Documentable) error {
	switch d := doc.(type) {
	case *document.Estimate:
		return r.UpdateEstimate(ctx, d)
	case *document.Mandate:
		return r.UpdateMandate(ctx, d)
	case *document.Contract:
		return r.UpdateContract(ctx, d)
	case *document.Invoice:
		return r.UpdateInvoice(ctx, d)
	default:
		return fmt.Errorf("unsupported document type: %T", doc)
	}
}