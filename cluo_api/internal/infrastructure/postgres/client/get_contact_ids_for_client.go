package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) GetContactIDsForClient(ctx context.Context, clientID uuid.UUID) ([]uuid.UUID, error) {
	query := fmt.Sprintf(`
		SELECT id
		FROM %s.contacts
		WHERE client_id = $1
		ORDER BY created_at
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, clientID)
	if err != nil {
		return nil, errs.ClassifyPgError("get contact IDs for client", err)
	}
	defer rows.Close()

	var contactIDs []uuid.UUID

	for rows.Next() {
		var contactID uuid.UUID
		err := rows.Scan(&contactID)
		if err != nil {
			return nil, errs.ClassifyPgError("scan contact ID", err)
		}
		contactIDs = append(contactIDs, contactID)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate contact IDs rows", err)
	}

	return contactIDs, nil
}