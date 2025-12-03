package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) UpdateClient(ctx context.Context, client *client.ClientEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.clients SET
			name_encrypted = $2,
			name_hash = $3,
			type_encrypted = $4,
			type_hash = $5,
			contactids_encrypted = $6,
			dek_encrypted = $7,
			key_version = $8,
			metadata = $9
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		client.ID, client.NameEncrypted, client.NameHash, client.TypeEncrypted,
		client.TypeHash, client.ContactIDsEncrypted, client.DEKEncrypted,
		client.KeyVersion, client.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("update client", err)
	}

	// Check if any row was actually updated
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}