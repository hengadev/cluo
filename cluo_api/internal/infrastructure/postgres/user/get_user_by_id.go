package userRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/user"
)

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.UserEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, email_hash, email_encrypted, password_hash_secure, role_encrypted,
			created_at_encrypted, dek_encrypted, key_version, metadata
		FROM %s.users WHERE id = $1
	`, r.schema)

	userEncx := &user.UserEncx{}

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&userEncx.ID, &userEncx.EmailHash, &userEncx.EmailEncrypted, &userEncx.PasswordHashSecure,
		&userEncx.RoleEncrypted, &userEncx.CreatedAtEncrypted,
		&userEncx.DEKEncrypted, &userEncx.KeyVersion, &userEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get user by id", err)
	}
	return userEncx, nil
}
