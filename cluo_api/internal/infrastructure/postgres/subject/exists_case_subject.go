package subjectRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) ExistsCaseSubject(ctx context.Context, id uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.case_subjects
			WHERE id = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check case subject existence: %w", err)
	}

	return exists, nil
}
