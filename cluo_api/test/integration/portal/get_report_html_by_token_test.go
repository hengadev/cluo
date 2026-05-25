package portal_test

import (
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReportHTMLByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("valid token with rapport returns 200 and non-empty HTML body", func(t *testing.T) {
		defer clearPortalTables(t, ctx)

		caseID := setupCaseWithRapport(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/html", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, body, "HTML body should not be empty")
	})

	t.Run("valid token with no rapport returns 404", func(t *testing.T) {
		defer clearPortalTables(t, ctx)

		c := clientHelpers.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

		testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
		require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    testCase.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/html", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		defer clearPortalTables(t, ctx)

		caseID := setupCaseWithRapport(t, ctx)

		rawToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    caseID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/html", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("revoked token returns 401", func(t *testing.T) {
		defer clearPortalTables(t, ctx)

		caseID := setupCaseWithRapport(t, ctx)

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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/html", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
