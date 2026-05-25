package portal_test

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/hengadev/cluo_api/internal/domain/token"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	investigationHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	rapportHelpers "github.com/hengadev/cluo_api/test/helpers/rapport"
	mediaHelpers "github.com/hengadev/cluo_api/test/helpers/media"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const archiveTipTapJSON = `{"type":"doc","content":[{"type":"heading","attrs":{"level":1},"content":[{"type":"text","text":"Rapport pour archive"}]},{"type":"paragraph","content":[{"type":"text","text":"Contenu du rapport."}]}]}`

// setupCaseWithArchiveData creates a client, a case, a rapport, and published media.
// Returns the case ID.
func setupCaseWithArchiveData(t *testing.T, ctx context.Context) uuid.UUID {
	t.Helper()

	c := clientHelpers.NewTestClient(t)
	clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
	require.NoError(t, err)
	require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

	testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
	require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))

	// Insert a rapport
	r := &rapport.Rapport{
		ID:        uuid.New(),
		CaseID:    testCase.ID,
		Content:   []byte(archiveTipTapJSON),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	rEncx, err := rapport.ProcessRapportEncx(ctx, crypto, r)
	require.NoError(t, err)
	require.NoError(t, rapportHelpers.InsertRapportEncx(t, ctx, testPool, rEncx))

	// Insert published media
	mediaFile := &domainMedia.MediaFile{
		ID:          uuid.New(),
		CaseID:      testCase.ID,
		URL:         "https://storage.example.com/test-photo.jpg",
		Type:        domainMedia.MediaTypeImage,
		MimeType:    "image/jpeg",
		FileName:    "test-photo.jpg",
		FileSize:    1024,
		Caption:     "Test photo",
		IsPublished: true,
		CreatedAt:   time.Now(),
	}
	mediaEncx, err := domainMedia.ProcessMediaFileEncx(ctx, crypto, mediaFile)
	require.NoError(t, err)
	mediaHelpers.InsertMediaEncx(t, ctx, testPool, mediaEncx)

	return testCase.ID
}

func clearArchiveTables(t *testing.T, ctx context.Context) {
	t.Helper()
	mediaHelpers.ClearMediaTable(t, ctx, testPool)
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE cases.rapports RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE cases.case_access_tokens RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	investigationHelpers.ClearCasesTable(t, ctx, testPool)
	clientHelpers.ClearClientsTable(t, ctx, testPool)
}

func TestGetFullArchiveByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("valid token returns 200 with zip content type", func(t *testing.T) {
		defer clearArchiveTables(t, ctx)

		caseID := setupCaseWithArchiveData(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/archive", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/zip", resp.Header.Get("Content-Type"))
		assert.Contains(t, resp.Header.Get("Content-Disposition"), "dossier-")
		assert.Contains(t, resp.Header.Get("Content-Disposition"), ".zip")

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, body)

		// Verify it's a valid zip
		reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		require.NoError(t, err)
		assert.NotEmpty(t, reader.File, "archive should contain entries")
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		defer clearArchiveTables(t, ctx)

		caseID := setupCaseWithArchiveData(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/archive", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("revoked token returns 401", func(t *testing.T) {
		defer clearArchiveTables(t, ctx)

		caseID := setupCaseWithArchiveData(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		revokedAt := time.Now().Add(-1 * time.Hour)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			RevokedAt: &revokedAt,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/archive", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestGetMediaArchiveByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("valid token returns 200 with zip content type", func(t *testing.T) {
		defer clearArchiveTables(t, ctx)

		caseID := setupCaseWithArchiveData(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media/archive", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/zip", resp.Header.Get("Content-Type"))
		assert.Contains(t, resp.Header.Get("Content-Disposition"), "medias.zip")

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, body)

		// Verify it's a valid zip with entries (testArchiveAdapter returns fake S3 content)
		reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		require.NoError(t, err)
		assert.NotEmpty(t, reader.File, "media archive should contain entries")
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		defer clearArchiveTables(t, ctx)

		caseID := setupCaseWithArchiveData(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media/archive", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("revoked token returns 401", func(t *testing.T) {
		defer clearArchiveTables(t, ctx)

		caseID := setupCaseWithArchiveData(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		revokedAt := time.Now().Add(-1 * time.Hour)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			RevokedAt: &revokedAt,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media/archive", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
