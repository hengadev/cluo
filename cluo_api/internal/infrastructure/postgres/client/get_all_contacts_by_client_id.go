package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (r *Repository) GetAllContactsByClientID(ctx context.Context, clientID uuid.UUID) ([]*client.ContactEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, client_id, firstname_encrypted, lastname_encrypted, email_encrypted,
			email_hash, phone_encrypted, position_encrypted, dek_encrypted, key_version
		FROM %s.contacts
		WHERE client_id = $1
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, errs.ClassifyPgError("get all contacts by client ID", err)
	}
	defer rows.Close()

	var contactsEncx []*client.ContactEncx

	for rows.Next() {
		contactEncx := &client.ContactEncx{}

		err := rows.Scan(
			&contactEncx.ID, &contactEncx.ClientID, &contactEncx.FirstnameEncrypted, &contactEncx.LastnameEncrypted, &contactEncx.EmailEncrypted,
			&contactEncx.EmailHash, &contactEncx.PhoneEncrypted,
			&contactEncx.PositionEncrypted, &contactEncx.DEKEncrypted,
			&contactEncx.KeyVersion,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan contact by client ID", err)
		}

		contactsEncx = append(contactsEncx, contactEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate contracts by client ID rows", err)
	}

	if len(contactsEncx) == 0 {
		return []*client.ContactEncx{}, nil
	}

	return contactsEncx, nil
}
