package subject

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domainsubject "github.com/hengadev/cluo_api/internal/domain/subject"
)

func (r *Repository) CreateCaseSubject(ctx context.Context, s *domainsubject.SubjectEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.case_subjects (
			id, created_at,
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
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		s.ID, s.CreatedAt,
		s.LastnameEncrypted, s.LastnameHash,
		s.FirstnameEncrypted, s.FirstnameHash,
		s.EmailEncrypted, s.EmailHash,
		s.PhoneEncrypted,
		s.CityEncrypted, s.CityHash,
		s.PostalCodeEncrypted, s.PostalCodeHash,
		s.Address1Encrypted, s.Address1Hash,
		s.Address2Encrypted, s.Address2Hash,
		s.OccupationEncrypted, s.OccupationHash,
		s.NotesEncrypted,
		s.DEKEncrypted, s.KeyVersion, s.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create case subject", err)
	}
	return nil
}
