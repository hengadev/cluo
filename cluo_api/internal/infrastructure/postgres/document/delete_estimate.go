package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// DeleteEstimate deletes an estimate by its ID.
func (r *Repository) DeleteEstimate(ctx context.Context, id string) error {
	query := `DELETE FROM estimates WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete estimate: %w", err)
	}
	return nil
}