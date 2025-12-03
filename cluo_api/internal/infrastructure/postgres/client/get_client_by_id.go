package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) GetClientByID(ctx context.Context, clientID uuid.UUID) (*client.ClientEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, created_at, name_encrypted, name_hash, type_encrypted, type_hash,
			contactids_encrypted, dek_encrypted, key_version, metadata
		FROM %s.clients WHERE id = $1
	`, r.schema)

	clientEncx := &client.ClientEncx{}

	err := r.pool.QueryRow(ctx, query, clientID).Scan(
		&clientEncx.ID, &clientEncx.CreatedAt, &clientEncx.NameEncrypted, &clientEncx.NameHash,
		&clientEncx.TypeEncrypted, &clientEncx.TypeHash, &clientEncx.ContactIDsEncrypted,
		&clientEncx.DEKEncrypted, &clientEncx.KeyVersion, &clientEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get client by id", err)
	}
	return clientEncx, nil
}