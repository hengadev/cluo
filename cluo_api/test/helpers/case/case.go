package caseHelpers

import (
	"testing"
	"time"

	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// NewTestCase creates a Case domain object with basic test data (plaintext fields only)
func NewTestCase(t *testing.T) *caseDomain.Case {
	t.Helper()

	contactID := uuid.New()
	return &caseDomain.Case{
		ID:                uuid.New(),
		Title:             "Test Case Title",
		Description:       "Test case description for unit testing",
		ClientID:          uuid.New(),
		AssignedContactID: &contactID,
		Status:            caseDomain.CaseStatusDraft,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// NewTestCaseEncx creates a mock CaseEncx domain object with basic test data (plaintext fields only)
func NewTestCaseEncx(t *testing.T) *caseDomain.CaseEncx {
	t.Helper()
	contactID := uuid.New()
	return &caseDomain.CaseEncx{
		ID:                   uuid.New(),
		ClientID:             uuid.New(),
		AssignedContactID:    &contactID,
		CreatedAt:            time.Now(),
		TitleEncrypted:       []byte("title_encrypted"),
		DescriptionEncrypted: []byte("description_encrypted"),
		StatusEncrypted:      []byte("status_encrypted"),
		UpdatedAtEncrypted:   []byte("updatedat_encrypted"),
		DEKEncrypted:         []byte("dek_encrypted"),
		KeyVersion:           1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

// NewTestCaseEncxWithClientID creates a mock CaseEncx with a specific client ID
func NewTestCaseEncxWithClientID(t *testing.T, clientID uuid.UUID) *caseDomain.CaseEncx {
	t.Helper()
	contactID := uuid.New()
	return &caseDomain.CaseEncx{
		ID:                   uuid.New(),
		ClientID:             clientID,
		AssignedContactID:    &contactID,
		CreatedAt:            time.Now(),
		TitleEncrypted:       []byte("title_encrypted"),
		DescriptionEncrypted: []byte("description_encrypted"),
		StatusEncrypted:      []byte("status_encrypted"),
		UpdatedAtEncrypted:   []byte("updatedat_encrypted"),
		DEKEncrypted:         []byte("dek_encrypted"),
		KeyVersion:           1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

// NewTestCaseEncxWithTimestamp creates a mock CaseEncx with a specific timestamp for ordering tests
func NewTestCaseEncxWithTimestamp(t *testing.T, clientID uuid.UUID, timestampOffset int) *caseDomain.CaseEncx {
	t.Helper()
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	contactID := uuid.New()
	return &caseDomain.CaseEncx{
		ID:                   uuid.New(),
		ClientID:             clientID,
		AssignedContactID:    &contactID,
		CreatedAt:            baseTime.Add(time.Duration(timestampOffset) * time.Hour),
		TitleEncrypted:       []byte("title_encrypted"),
		DescriptionEncrypted: []byte("description_encrypted"),
		StatusEncrypted:      []byte("status_encrypted"),
		UpdatedAtEncrypted:   []byte("updatedat_encrypted"),
		DEKEncrypted:         []byte("dek_encrypted"),
		KeyVersion:           1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

