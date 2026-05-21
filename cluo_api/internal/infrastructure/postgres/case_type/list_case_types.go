package caseTypeRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (r *Repository) ListCaseTypes(ctx context.Context) ([]*casetype.CaseType, error) {
	query := fmt.Sprintf(`
		SELECT id, name, created_at, updated_at
		FROM %s.case_types
		ORDER BY name ASC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("list case types", err)
	}
	defer rows.Close()

	var caseTypes []*casetype.CaseType
	for rows.Next() {
		ct := &casetype.CaseType{}
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.CreatedAt, &ct.UpdatedAt); err != nil {
			return nil, errs.ClassifyPgError("scan case type row", err)
		}
		caseTypes = append(caseTypes, ct)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate case type rows", err)
	}

	return caseTypes, nil
}
