package mediaHelpers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/require"
)

// NewTestMedia creates a MediaFile domain object with basic test data
func NewTestMedia(t *testing.T) *domain.MediaFile {
	t.Helper()

	return &domain.MediaFile{
		ID:          uuid.New(),
		CaseID:      uuid.New(),
		URL:         "https://storage.example.com/test-file.jpg",
		Type:        domain.MediaTypeImage,
		MimeType:    "image/jpeg",
		FileName:    "test-file.jpg",
		FileSize:    1024000,
		Caption:     "Test image caption",
		IsPublished: false,
		CreatedAt:   time.Now(),
	}
}

// NewTestMediaWithCaseID creates media with specific case ID
func NewTestMediaWithCaseID(t *testing.T, caseID uuid.UUID) *domain.MediaFile {
	t.Helper()

	media := NewTestMedia(t)
	media.CaseID = caseID
	return media
}

// NewTestMediaEncx creates a mock MediaFileEncx object
func NewTestMediaEncx(t *testing.T) *domain.MediaFileEncx {
	t.Helper()

	return &domain.MediaFileEncx{
		ID:                  uuid.New(),
		CaseID:              uuid.New(),
		FileSize:            1024000,
		IsPublished:         false,
		CreatedAt:           time.Now(),
		URLEncrypted:        []byte("url_encrypted"),
		TypeEncrypted:       []byte("type_encrypted"),
		MimeTypeEncrypted:   []byte("mimetype_encrypted"),
		FileNameEncrypted:   []byte("filename_encrypted"),
		CaptionEncrypted:    []byte("caption_encrypted"),
		DEKEncrypted:        []byte("dek_encrypted"),
		KeyVersion:          1,
		Metadata: encx.EncryptionMetadata{
			KEKAlias:         "test-kek-alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "1.0.0",
		},
	}
}

// NewTestMediaEncxWithCaseID creates MediaEncx with specific case ID
func NewTestMediaEncxWithCaseID(t *testing.T, caseID uuid.UUID) *domain.MediaFileEncx {
	t.Helper()

	mediaEncx := NewTestMediaEncx(t)
	mediaEncx.CaseID = caseID
	return mediaEncx
}

// CreateEncryptedTestMedia creates a properly encrypted MediaFileEncx for testing
// This should be used in tests that need to insert media directly into the database
func CreateEncryptedTestMedia(t *testing.T, ctx context.Context, crypto encx.CryptoService, caseID uuid.UUID) *domain.MediaFileEncx {
	t.Helper()

	// Create plain media
	media := NewTestMediaWithCaseID(t, caseID)

	// Encrypt it properly using ProcessMediaFileEncx
	mediaEncx, err := domain.ProcessMediaFileEncx(ctx, crypto, media)
	require.NoError(t, err)

	return mediaEncx
}
