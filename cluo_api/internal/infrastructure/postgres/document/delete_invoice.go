package documentRepository

import (
	"context"
	"fmt"
)

// DeleteInvoice deletes an invoice by its ID.
func (r *Repository) DeleteInvoice(ctx context.Context, id string) error {
	query := `DELETE FROM invoices WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}
	return nil
}
