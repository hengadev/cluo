package portal_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
	"github.com/hengadev/cluo_api/internal/domain/token"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	investigationHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	rapportHelpers "github.com/hengadev/cluo_api/test/helpers/rapport"
	"github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sampleTipTapJSON is a minimal valid TipTap document for test rapport content.
const sampleTipTapJSON = `{"type":"doc","content":[{"type":"heading","attrs":{"level":1},"content":[{"type":"text","text":"Rapport de test"}]},{"type":"paragraph","content":[{"type":"text","text":"Contenu du rapport."}]}]}`

func clearPortalTables(t *testing.T, ctx context.Context) {
	t.Helper()
	_, err := testPool.Exec(ctx, "TRUNCATE TABLE cases.case_access_tokens RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE TABLE cases.rapports RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	investigationHelpers.ClearCasesTable(t, ctx, testPool)
	clientHelpers.ClearClientsTable(t, ctx, testPool)
}

// insertTestToken inserts a token record directly into the database.
func insertTestToken(t *testing.T, ctx context.Context, tk *token.Token) {
	t.Helper()
	_, err := testPool.Exec(ctx,
		`INSERT INTO cases.case_access_tokens (id, case_id, token_hash, expires_at, revoked_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		tk.ID, tk.CaseID, tk.TokenHash, tk.ExpiresAt, tk.RevokedAt, tk.CreatedAt,
	)
	require.NoError(t, err)
}

// setupCaseWithRapport creates a client, a case, and an encrypted rapport. Returns the case ID.
func setupCaseWithRapport(t *testing.T, ctx context.Context) uuid.UUID {
	t.Helper()

	c := clientHelpers.NewTestClient(t)
	clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
	require.NoError(t, err)
	require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

	testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
	require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))

	r := &rapport.Rapport{
		ID:        uuid.New(),
		CaseID:    testCase.ID,
		Content:   []byte(sampleTipTapJSON),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	rEncx, err := rapport.ProcessRapportEncx(ctx, crypto, r)
	require.NoError(t, err)
	require.NoError(t, rapportHelpers.InsertRapportEncx(t, ctx, testPool, rEncx))

	return testCase.ID
}

func TestGetReportPDFByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("valid token with rapport returns 200 and PDF", func(t *testing.T) {
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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/pdf", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/pdf", resp.Header.Get("Content-Type"))
		assert.NotEmpty(t, resp.Header.Get("Content-Length"))
		assert.Contains(t, resp.Header.Get("Content-Disposition"), "rapport.pdf")
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
			ExpiresAt: time.Now().Add(-1 * time.Hour), // already expired
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/pdf", testServerURL, rawToken))
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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/pdf", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("valid token with no rapport returns 404", func(t *testing.T) {
		defer clearPortalTables(t, ctx)

		// Insert a case with no rapport
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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/report/pdf", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
