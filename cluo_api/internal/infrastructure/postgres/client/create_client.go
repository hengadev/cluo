package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) CreateClient(ctx context.Context, clientEncx *client.ClientEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.clients (
			id, created_at, name_encrypted, name_hash, type_encrypted, type_hash,
			dek_encrypted, key_version, metadata
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		clientEncx.ID, clientEncx.CreatedAt, clientEncx.NameEncrypted, clientEncx.NameHash,
		clientEncx.TypeEncrypted, clientEncx.TypeHash,
		clientEncx.DEKEncrypted, clientEncx.KeyVersion, clientEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create client", err)
	}

	return nil
}
