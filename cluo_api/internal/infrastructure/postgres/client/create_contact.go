package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) CreateContact(ctx context.Context, contactEncx *client.ContactEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.contacts (
			id, client_id_hash, client_id_encrypted, lastname_encrypted, firstname_encrypted, email_hash, email_encrypted, 
			phone_encrypted, position_encrypted, dek_encrypted, key_version
		) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		contactEncx.ID, contactEncx.ClientIDHash, contactEncx.ClientIDEncrypted, contactEncx.LastnameEncrypted, contactEncx.FirstnameEncrypted, contactEncx.EmailHash,
		contactEncx.EmailEncrypted, contactEncx.PhoneEncrypted, contactEncx.PositionEncrypted, contactEncx.DEKEncrypted, contactEncx.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("create contact", err)
	}

	return nil
}
