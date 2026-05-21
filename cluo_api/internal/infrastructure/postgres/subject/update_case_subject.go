package subject

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domainsubject "github.com/hengadev/cluo_api/internal/domain/subject"
)

func (r *Repository) UpdateCaseSubject(ctx context.Context, s *domainsubject.SubjectEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.case_subjects SET
			lastname_encrypted=$1,  lastname_hash=$2,
			firstname_encrypted=$3, firstname_hash=$4,
			email_encrypted=$5,     email_hash=$6,
			phone_encrypted=$7,
			city_encrypted=$8,      city_hash=$9,
			postal_code_encrypted=$10, postal_code_hash=$11,
			address1_encrypted=$12, address1_hash=$13,
			address2_encrypted=$14, address2_hash=$15,
			occupation_encrypted=$16, occupation_hash=$17,
			notes_encrypted=$18,
			dek_encrypted=$19, key_version=$20, metadata=$21
		WHERE id=$22`, r.schema)

	result, err := r.pool.Exec(ctx, query,
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
		s.ID,
	)
	if err != nil {
		return errs.ClassifyPgError("update case subject", err)
	}
	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("case subject %s not found", s.ID), "case_subject")
	}
	return nil
}
