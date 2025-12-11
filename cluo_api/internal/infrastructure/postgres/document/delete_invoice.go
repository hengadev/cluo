package documentRepository

import (
	"context"
	"fmt"
)

	query := `DELETE FROM invoices WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}
	return nil
}
