package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// Delete deletes a document by its ID and type.
func (r *Repository) Delete(ctx context.Context, id string, docType document.DocumentType) error {
	switch docType {
	case document.DocumentTypeEstimate:
		return r.DeleteEstimate(ctx, id)
	case document.DocumentTypeMandate:
		return r.DeleteMandate(ctx, id)
	case document.DocumentTypeContract:
		return r.DeleteContract(ctx, id)
	case document.DocumentTypeInvoice:
		return r.DeleteInvoice(ctx, id)
	default:
		return fmt.Errorf("unsupported document type: %s", docType)
	}
}
