package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetByID retrieves a document by its ID and type.
func (r *Repository) GetByID(ctx context.Context, id string, docType document.DocumentType) (document.Documentable, error) {
	switch docType {
	case document.DocumentTypeEstimate:
		return r.GetEstimateByID(ctx, id)
	case document.DocumentTypeMandate:
		return r.GetMandateByID(ctx, id)
	case document.DocumentTypeContract:
		return r.GetContractByID(ctx, id)
	case document.DocumentTypeInvoice:
		return r.GetInvoiceByID(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported document type: %s", docType)
	}
}
