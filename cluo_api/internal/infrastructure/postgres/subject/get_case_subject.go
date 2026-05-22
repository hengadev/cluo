package subjectRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	domainsubject "github.com/hengadev/cluo_api/internal/domain/subject"
)

func (r *Repository) GetCaseSubjectByID(ctx context.Context, id uuid.UUID) (*domainsubject.SubjectEncx, error) {
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
		FROM %s.case_subjects WHERE id = $1`, r.schema)

	s := &domainsubject.SubjectEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
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
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get case subject by id", err)
	}
	return s, nil
}
