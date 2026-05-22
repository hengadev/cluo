package subjectRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domainsubject "github.com/hengadev/cluo_api/internal/domain/subject"
)

func (r *Repository) ListCaseSubjects(ctx context.Context, page, pageSize int) ([]*domainsubject.SubjectEncx, int, error) {
	offset := (page - 1) * pageSize

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.case_subjects`, r.schema)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errs.ClassifyPgError("count case subjects", err)
	}
	if total == 0 {
		return []*domainsubject.SubjectEncx{}, 0, nil
	}

	query := fmt.Sprintf(`
		SELECT
			id, created_at AS createdat,
			lastname_encrypted, lastname_hash,
			firstname_encrypted, firstname_hash,
			email_encrypted, email_hash,
			phone_encrypted,
			city_encrypted, city_hash,
			postal_code_encrypted, postal_code_hash,
			address1_encrypted, address1_hash,
			address2_encrypted, address2_hash,
			occupation_encrypted, occupation_hash,
			notes_encrypted,
			dek_encrypted, key_version, metadata
		FROM %s.case_subjects
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, r.schema)

	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errs.ClassifyPgError("list case subjects", err)
	}
	defer rows.Close()

	var subjects []*domainsubject.SubjectEncx
	for rows.Next() {
		s := &domainsubject.SubjectEncx{}
		if err := rows.Scan(
			&s.ID, &s.CreatedAt,
			&s.LastnameEncrypted, &s.LastnameHash,
			&s.FirstnameEncrypted, &s.FirstnameHash,
			&s.EmailEncrypted, &s.EmailHash,
			&s.PhoneEncrypted,
			&s.CityEncrypted, &s.CityHash,
			&s.PostalCodeEncrypted, &s.PostalCodeHash,
			&s.Address1Encrypted, &s.Address1Hash,
			&s.Address2Encrypted, &s.Address2Hash,
			&s.OccupationEncrypted, &s.OccupationHash,
			&s.NotesEncrypted,
			&s.DEKEncrypted, &s.KeyVersion, &s.Metadata,
		); err != nil {
			return nil, 0, errs.ClassifyPgError("scan case subject", err)
		}
		subjects = append(subjects, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errs.ClassifyPgError("iterate case subjects", err)
	}

	return subjects, total, nil
}
