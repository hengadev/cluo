package investigationHelpers

import (
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/investigation"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// NewTestCase creates a Case domain object with basic test data (plaintext fields only)
func NewTestCase(t *testing.T) *investigation.Investigation {
	t.Helper()

	contactID := uuid.New()
	externalRef := "EXT-REF-123"
	latitude := "48.8566"
	longitude := "2.3522"
	return &investigation.Investigation{
		ID:                uuid.New(),
		Title:             "Test Case Title",
		Description:       "Test case description for unit testing",
		ClientID:          uuid.New(),
		AssignedContactID: &contactID,
		CaseSubjectID:     nil,
		ExternalReference: &externalRef,
		CaseTypeID:        nil,
		Status:            investigation.StatusInProgress,
		Placename:         "Test Location",
		Address1:          "123 Test Street",
		Address2:          "Apt 4B",
		City:              "Test City",
		PostalCode:        "12345",
		Country:           "Test Country",
		Latitude:          &latitude,
		Longitude:         &longitude,
		LocationType:      "residence",
		LocationNotes:     "Test location notes",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// NewTestCaseEncx creates a mock CaseEncx domain object with basic test data (plaintext fields only)
func NewTestCaseEncx(t *testing.T) *investigation.InvestigationEncx {
	t.Helper()
	contactID := uuid.New()
	return &investigation.InvestigationEncx{
		ID:                         uuid.New(),
		ClientID:                   uuid.New(),
		AssignedContactID:          &contactID,
		CaseSubjectID:              nil,
		CaseTypeID:                 nil,
		CreatedAt:                  time.Now(),
		TitleEncrypted:             []byte("title_encrypted"),
		DescriptionEncrypted:       []byte("description_encrypted"),
		ExternalReferenceEncrypted: []byte("external_ref_encrypted"),
		ExternalReferenceHash:      "external_ref_hash",
		StatusEncrypted:            []byte("status_encrypted"),
		PlacenameEncrypted:         []byte("placename_encrypted"),
		PlacenameHash:              "placename_hash",
		Address1Encrypted:          []byte("address1_encrypted"),
		Address1Hash:               "address1_hash",
		Address2Encrypted:          []byte("address2_encrypted"),
		Address2Hash:               "address2_hash",
		CityEncrypted:              []byte("city_encrypted"),
		CityHash:                   "city_hash",
		PostalCodeEncrypted:        []byte("postal_code_encrypted"),
		PostalCodeHash:             "postal_code_hash",
		CountryEncrypted:           []byte("country_encrypted"),
		CountryHash:                "country_hash",
		LatitudeEncrypted:          []byte("latitude_encrypted"),
		LatitudeHash:               "latitude_hash",
		LongitudeEncrypted:         []byte("longitude_encrypted"),
		LongitudeHash:              "longitude_hash",
		LocationTypeEncrypted:      []byte("location_type_encrypted"),
		LocationTypeHash:           "location_type_hash",
		LocationNotesEncrypted:     []byte("location_notes_encrypted"),
		LocationNotesHash:          "location_notes_hash",
		UpdatedAtEncrypted:         []byte("updatedat_encrypted"),
		DEKEncrypted:               []byte("dek_encrypted"),
		KeyVersion:                 1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

// NewTestCaseEncxWithClientID creates a mock CaseEncx with a specific client ID
func NewTestCaseEncxWithClientID(t *testing.T, clientID uuid.UUID) *investigation.InvestigationEncx {
	t.Helper()
	return &investigation.InvestigationEncx{
		ID:                         uuid.New(),
		ClientID:                   clientID,
		AssignedContactID:          nil, // No contact assigned
		CaseSubjectID:              nil,
		CaseTypeID:                 nil,
		CreatedAt:                  time.Now(),
		TitleEncrypted:             []byte("title_encrypted"),
		DescriptionEncrypted:       []byte("description_encrypted"),
		ExternalReferenceEncrypted: nil,
		ExternalReferenceHash:      "",
		StatusEncrypted:            []byte("status_encrypted"),
		PlacenameEncrypted:         []byte("placename_encrypted"),
		PlacenameHash:              "placename_hash",
		Address1Encrypted:          []byte("address1_encrypted"),
		Address1Hash:               "address1_hash",
		Address2Encrypted:          []byte("address2_encrypted"),
		Address2Hash:               "address2_hash",
		CityEncrypted:              []byte("city_encrypted"),
		CityHash:                   "city_hash",
		PostalCodeEncrypted:        []byte("postal_code_encrypted"),
		PostalCodeHash:             "postal_code_hash",
		CountryEncrypted:           []byte("country_encrypted"),
		CountryHash:                "country_hash",
		LatitudeEncrypted:          []byte("latitude_encrypted"),
		LatitudeHash:               "latitude_hash",
		LongitudeEncrypted:         []byte("longitude_encrypted"),
		LongitudeHash:              "longitude_hash",
		LocationTypeEncrypted:      []byte("location_type_encrypted"),
		LocationTypeHash:           "location_type_hash",
		LocationNotesEncrypted:     []byte("location_notes_encrypted"),
		LocationNotesHash:          "location_notes_hash",
		UpdatedAtEncrypted:         []byte("updatedat_encrypted"),
		DEKEncrypted:               []byte("dek_encrypted"),
		KeyVersion:                 1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

// NewTestCaseEncxWithTimestamp creates a mock CaseEncx with a specific timestamp for ordering tests
func NewTestCaseEncxWithTimestamp(t *testing.T, clientID uuid.UUID, timestampOffset int) *investigation.InvestigationEncx {
	t.Helper()
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	contactID := uuid.New()
	return &investigation.InvestigationEncx{
		ID:                         uuid.New(),
		ClientID:                   clientID,
		AssignedContactID:          &contactID,
		CaseSubjectID:              nil,
		CaseTypeID:                 nil,
		CreatedAt:                  baseTime.Add(time.Duration(timestampOffset) * time.Hour),
		TitleEncrypted:             []byte("title_encrypted"),
		DescriptionEncrypted:       []byte("description_encrypted"),
		ExternalReferenceEncrypted: []byte("external_ref_encrypted"),
		ExternalReferenceHash:      "external_ref_hash",
		StatusEncrypted:            []byte("status_encrypted"),
		PlacenameEncrypted:         []byte("placename_encrypted"),
		PlacenameHash:              "placename_hash",
		Address1Encrypted:          []byte("address1_encrypted"),
		Address1Hash:               "address1_hash",
		Address2Encrypted:          []byte("address2_encrypted"),
		Address2Hash:               "address2_hash",
		CityEncrypted:              []byte("city_encrypted"),
		CityHash:                   "city_hash",
		PostalCodeEncrypted:        []byte("postal_code_encrypted"),
		PostalCodeHash:             "postal_code_hash",
		CountryEncrypted:           []byte("country_encrypted"),
		CountryHash:                "country_hash",
		LatitudeEncrypted:          []byte("latitude_encrypted"),
		LatitudeHash:               "latitude_hash",
		LongitudeEncrypted:         []byte("longitude_encrypted"),
		LongitudeHash:              "longitude_hash",
		LocationTypeEncrypted:      []byte("location_type_encrypted"),
		LocationTypeHash:           "location_type_hash",
		LocationNotesEncrypted:     []byte("location_notes_encrypted"),
		LocationNotesHash:          "location_notes_hash",
		UpdatedAtEncrypted:         []byte("updatedat_encrypted"),
		DEKEncrypted:               []byte("dek_encrypted"),
		KeyVersion:                 1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}
