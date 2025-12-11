package documentRepository

import (
	"context"
	"fmt"
)

// DeleteContract deletes a contract by its ID.
func (r *Repository) DeleteContract(ctx context.Context, id string) error {
	query := `DELETE FROM contracts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}
	return nil
}
