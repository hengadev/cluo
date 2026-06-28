package userRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) UpdateUserName(ctx context.Context, userID uuid.UUID, nameEncrypted []byte) error {
	query := fmt.Sprintf(`
		UPDATE %s.users SET name_encrypted = $1 WHERE id = $2
	`, r.schema)

	tag, err := r.pool.Exec(ctx, query, nameEncrypted, userID)
	if err != nil {
		return errs.ClassifyPgError("update user name", err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrRepositoryNotFound
	}
	return nil
}
