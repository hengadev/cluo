package subjectHelpers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	subjectDomain "github.com/hengadev/cluo_api/internal/domain/subject"
)

// NewTestSubjectEncx creates a mock SubjectEncx with basic test data.
func NewTestSubjectEncx(t *testing.T) *subjectDomain.SubjectEncx {
	t.Helper()
	return &subjectDomain.SubjectEncx{
		ID:                      uuid.New(),
		LastnameEncrypted:       []byte("lastname_encrypted"),
		LastnameHash:            "lastname_hash",
		FirstnameEncrypted:      []byte("firstname_encrypted"),
		FirstnameHash:           "firstname_hash",
		EmailEncrypted:          []byte("email_encrypted"),
		EmailHash:               "email_hash",
		PhoneEncrypted:          []byte("phone_encrypted"),
		CityEncrypted:           []byte("city_encrypted"),
		CityHash:                "city_hash",
		PostalCodeEncrypted:     []byte("postal_code_encrypted"),
		PostalCodeHash:          "postal_code_hash",
		Address1Encrypted:       []byte("address1_encrypted"),
		Address1Hash:            "address1_hash",
		Address2Encrypted:       []byte("address2_encrypted"),
		Address2Hash:            "address2_hash",
		OccupationEncrypted:     []byte("occupation_encrypted"),
		OccupationHash:          "occupation_hash",
		NotesEncrypted:          []byte("notes_encrypted"),
		DEKEncrypted:            []byte("dek_encrypted"),
		KeyVersion:              1,
		Metadata:                encx.EncryptionMetadata{},
		CreatedAt:               time.Now(),
	}
}
