package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) GetContactByID(ctx context.Context, contactID uuid.UUID) (*client.ContactEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, client_id, firstname_encrypted, lastname_encrypted, email_encrypted, email_hash, 
			phone_encrypted, position_encrypted, dek_encrypted, key_version, created_at
		FROM %s.contacts WHERE id = $1
	`, r.schema)

	contactEncx := &client.ContactEncx{}

	err := r.pool.QueryRow(ctx, query, contactID).Scan(
		&contactEncx.ID, &contactEncx.ClientID, &contactEncx.FirstnameEncrypted, &contactEncx.LastnameEncrypted, &contactEncx.EmailEncrypted,
		&contactEncx.EmailHash, &contactEncx.PhoneEncrypted,
		&contactEncx.PositionEncrypted, &contactEncx.DEKEncrypted,
		&contactEncx.KeyVersion, &contactEncx.CreatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get contact by id", err)
	}
	return contactEncx, nil
}
