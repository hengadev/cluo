package portal_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/hengadev/cluo_api/internal/domain/token"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	investigationHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	mediaHelpers "github.com/hengadev/cluo_api/test/helpers/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func insertPublishedMedia(t *testing.T, ctx context.Context, caseID uuid.UUID, mediaType domainMedia.MediaType, filename string) {
	t.Helper()
	m := &domainMedia.MediaFile{
		ID:          uuid.New(),
		CaseID:      caseID,
		URL:         fmt.Sprintf("https://storage.example.com/%s", filename),
		Type:        mediaType,
		MimeType:    "application/octet-stream",
		FileName:    filename,
		FileSize:    1024,
		IsPublished: true,
		CreatedAt:   time.Now(),
	}
	encx, err := domainMedia.ProcessMediaFileEncx(ctx, crypto, m)
	require.NoError(t, err)
	mediaHelpers.InsertMediaEncx(t, ctx, testPool, encx)
}

func setupCaseForMediaTest(t *testing.T, ctx context.Context) uuid.UUID {
	t.Helper()
	c := clientHelpers.NewTestClient(t)
	clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
	require.NoError(t, err)
	require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

	testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
	require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))
	return testCase.ID
}

func clearMediaByTokenTables(t *testing.T, ctx context.Context) {
	t.Helper()
	mediaHelpers.ClearMediaTable(t, ctx, testPool)
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE cases.case_access_tokens RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	investigationHelpers.ClearCasesTable(t, ctx, testPool)
	clientHelpers.ClearClientsTable(t, ctx, testPool)
}

func TestGetMediaByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("audio files are excluded from response", func(t *testing.T) {
		defer clearMediaByTokenTables(t, ctx)

		caseID := setupCaseForMediaTest(t, ctx)
		insertPublishedMedia(t, ctx, caseID, domainMedia.MediaTypeImage, "photo.jpg")
		insertPublishedMedia(t, ctx, caseID, domainMedia.MediaTypeVideo, "clip.mp4")
		insertPublishedMedia(t, ctx, caseID, domainMedia.MediaTypeAudio, "recording.mp3")

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var media []domainMedia.MediaResponse
		require.NoError(t, json.Unmarshal(body, &media))

		assert.Len(t, media, 2, "should return image and video only")
		for _, m := range media {
			assert.NotEqual(t, "audio", m.Type, "audio must never appear in portal media response")
		}
	})

	t.Run("returns empty array when all published media is audio", func(t *testing.T) {
		defer clearMediaByTokenTables(t, ctx)

		caseID := setupCaseForMediaTest(t, ctx)
		insertPublishedMedia(t, ctx, caseID, domainMedia.MediaTypeAudio, "recording.mp3")

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var media []domainMedia.MediaResponse
		require.NoError(t, json.Unmarshal(body, &media))
		assert.Empty(t, media, "should return empty array, not error")
	})

	t.Run("images and videos with IsPublished=true are returned", func(t *testing.T) {
		defer clearMediaByTokenTables(t, ctx)

		caseID := setupCaseForMediaTest(t, ctx)
		insertPublishedMedia(t, ctx, caseID, domainMedia.MediaTypeImage, "a.jpg")
		insertPublishedMedia(t, ctx, caseID, domainMedia.MediaTypeVideo, "b.mp4")

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var media []domainMedia.MediaResponse
		require.NoError(t, json.Unmarshal(body, &media))
		assert.Len(t, media, 2)
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		defer clearMediaByTokenTables(t, ctx)

		caseID := setupCaseForMediaTest(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("revoked token returns 401", func(t *testing.T) {
		defer clearMediaByTokenTables(t, ctx)

		caseID := setupCaseForMediaTest(t, ctx)

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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/media", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
