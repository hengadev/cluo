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
			client_id = $2,
			lastname_encrypted = $3,
			firstname_encrypted = $4,
			email_hash = $5,
			email_encrypted = $6,
			phone_encrypted = $7,
			position_encrypted = $8,
			dek_encrypted = $9,
			key_version = $10
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		contact.ID, contact.ClientID, contact.LastnameEncrypted,
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
