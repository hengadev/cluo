package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) GetAllClients(ctx context.Context) ([]*client.ClientEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, created_at, name_encrypted, name_hash, type_encrypted, type_hash,
			contactids_encrypted, dek_encrypted, key_version, metadata
		FROM %s.clients
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all clients", err)
	}
	defer rows.Close()

	var clients []*client.ClientEncx

	for rows.Next() {
		clientEncx := &client.ClientEncx{}
		err := rows.Scan(
			&clientEncx.ID, &clientEncx.CreatedAt, &clientEncx.NameEncrypted, &clientEncx.NameHash,
			&clientEncx.TypeEncrypted, &clientEncx.TypeHash, &clientEncx.ContactIDsEncrypted,
			&clientEncx.DEKEncrypted, &clientEncx.KeyVersion, &clientEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan client row", err)
		}
		clients = append(clients, clientEncx)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate client rows", err)
	}

	return clients, nil
}