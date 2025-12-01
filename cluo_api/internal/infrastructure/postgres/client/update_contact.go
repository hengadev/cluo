package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) UpdateContact(ctx context.Context, contact *client.ContactEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.contacts SET
			client_id_hash = $2,
			client_id_encrypted = $3,
			lastname_encrypted = $4,
			firstname_encrypted = $5,
			email_hash = $6,
			email_encrypted = $7,
			phone_encrypted = $8,
			position_encrypted = $9,
			dek_encrypted = $10,
			key_version = $11
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		contact.ID, contact.ClientIDHash, contact.ClientIDEncrypted, contact.LastnameEncrypted,
		contact.FirstnameEncrypted, contact.EmailHash, contact.EmailEncrypted, contact.PhoneEncrypted,
		contact.PositionEncrypted, contact.DEKEncrypted, contact.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("update contact", err)
	}

	// Check if any row was actually updated
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
