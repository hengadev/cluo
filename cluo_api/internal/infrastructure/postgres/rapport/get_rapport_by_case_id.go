package rapportRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (r *Repository) GetRapportByCaseID(ctx context.Context, caseID uuid.UUID) (*rapport.RapportEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, case_id, content_encrypted, dek_encrypted, key_version, created_at, updated_at
		FROM %s.rapports
		WHERE case_id = $1
	`, r.schema)

	rEncx := &rapport.RapportEncx{}

	err := r.pool.QueryRow(ctx, query, caseID).Scan(
		&rEncx.ID,
		&rEncx.CaseID,
		&rEncx.ContentEncrypted,
		&rEncx.DEKEncrypted,
		&rEncx.KeyVersion,
		&rEncx.CreatedAt,
		&rEncx.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get rapport by case id", err)
	}

	return rEncx, nil
}
