package userRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/user"
)

func (r *Repository) CreateUser(ctx context.Context, userEncx *user.UserEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.users (
			id, email_hash, email_encrypted, password_hash_secure, role_encrypted,
			created_at_encrypted, dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		userEncx.ID, userEncx.EmailHash, userEncx.EmailEncrypted, userEncx.PasswordHashSecure,
		userEncx.RoleEncrypted, userEncx.CreatedAtEncrypted,
		userEncx.DEKEncrypted, userEncx.KeyVersion, userEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create user", err)
	}

	return nil
}
