package caseRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) ExistsCase(ctx context.Context, caseID uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s.cases WHERE id = $1)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, caseID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check case existence: %w", err)
	}

	return exists, nil
}
