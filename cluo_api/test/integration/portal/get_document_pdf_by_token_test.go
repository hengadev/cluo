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
	documentDomain "github.com/hengadev/cluo_api/internal/domain/document"
	"github.com/hengadev/cluo_api/internal/domain/token"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	documentHelpers "github.com/hengadev/cluo_api/test/helpers/document"
	investigationHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// insertDocumentForCase creates a client, a case, inserts a document of the
// given type for that case, and creates a valid access token.
// Returns the raw token string for HTTP requests.
func insertDocumentForCase(t *testing.T, ctx context.Context, docType documentDomain.DocumentType) string {
	t.Helper()

	c := clientHelpers.NewTestClient(t)
	clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
	require.NoError(t, err)
	require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

	testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
	require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))

	caseID := testCase.ID
	clientID := c.ID

	switch docType {
	case documentDomain.DocumentTypeEstimate:
		estimate := documentHelpers.NewTestEstimateWithCaseID(t, caseID, clientID)
		estimateEncx, err := documentDomain.ProcessEstimateEncx(ctx, crypto, estimate)
		require.NoError(t, err)
		require.NoError(t, documentHelpers.InsertEstimateWithCaseID(t, ctx, testPool, caseID, clientID, estimateEncx))

	case documentDomain.DocumentTypeMandate:
		mandate := documentHelpers.NewTestMandateWithCaseID(t, caseID, clientID)
		mandateEncx, err := documentDomain.ProcessMandateEncx(ctx, crypto, mandate)
		require.NoError(t, err)
		require.NoError(t, documentHelpers.InsertMandateWithCaseID(t, ctx, testPool, caseID, clientID, mandateEncx))

	case documentDomain.DocumentTypeContract:
		contract := documentHelpers.NewTestContractWithCaseID(t, caseID, clientID)
		contractEncx, err := documentDomain.ProcessContractEncx(ctx, crypto, contract)
		require.NoError(t, err)
		require.NoError(t, documentHelpers.InsertContractWithCaseID(t, ctx, testPool, caseID, clientID, contractEncx))

	case documentDomain.DocumentTypeInvoice:
		invoice := documentHelpers.NewTestInvoiceWithCaseID(t, caseID, clientID)
		invoiceEncx, err := documentDomain.ProcessInvoiceEncx(ctx, crypto, invoice)
		require.NoError(t, err)
		require.NoError(t, documentHelpers.InsertInvoiceWithCaseID(t, ctx, testPool, caseID, clientID, invoiceEncx))
	}

	rawToken, tokenHash, err := token.GenerateRawToken()
	require.NoError(t, err)
	insertTestToken(t, ctx, &token.Token{
		ID:        uuid.New(),
		CaseID:    caseID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	})

	return rawToken
}

func clearDocumentPortalTables(t *testing.T, ctx context.Context) {
	t.Helper()
	documentHelpers.ClearAllDocumentTables(t, ctx, testPool)
	clearPortalTables(t, ctx)
}

func TestGetDocumentPDFByToken(t *testing.T) {
	ctx := context.Background()

	for _, docType := range []string{"estimate", "mandate", "contract", "invoice"} {
		t.Run(fmt.Sprintf("valid token with %s returns 200 and PDF", docType), func(t *testing.T) {
			defer clearDocumentPortalTables(t, ctx)

			rawToken := insertDocumentForCase(t, ctx, toDocumentType(docType))

			resp, err := http.Get(fmt.Sprintf("%s/token/%s/documents/%s/pdf", testServerURL, rawToken, docType))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "application/pdf", resp.Header.Get("Content-Type"))
			assert.Contains(t, resp.Header.Get("Content-Disposition"), fmt.Sprintf(`attachment; filename="%s.pdf"`, docType))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, body)
			// Verify it starts with PDF magic bytes
			assert.True(t, len(body) >= 4 && string(body[:4]) == "%PDF", "Response should be a valid PDF")
		})
	}

	t.Run("invalid document type returns 400", func(t *testing.T) {
		defer clearDocumentPortalTables(t, ctx)

		// Create a valid case + token but with no documents
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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/documents/invalidtype/pdf", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("valid token but document not found returns 404", func(t *testing.T) {
		defer clearDocumentPortalTables(t, ctx)

		// Create a valid case + token but with no documents
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

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/documents/estimate/pdf", testServerURL, rawToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		defer clearDocumentPortalTables(t, ctx)

		rawToken := insertDocumentForCase(t, ctx, documentDomain.DocumentTypeEstimate)

		// Overwrite the token with an expired one
		c := clientHelpers.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

		testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
		require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))

		expiredToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    testCase.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-48 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/documents/estimate/pdf", testServerURL, expiredToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		_ = rawToken
	})

	t.Run("revoked token returns 401", func(t *testing.T) {
		defer clearDocumentPortalTables(t, ctx)

		rawToken := insertDocumentForCase(t, ctx, documentDomain.DocumentTypeEstimate)

		// Overwrite the token with a revoked one
		c := clientHelpers.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		require.NoError(t, clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx))

		testCase := investigationHelpers.NewTestCaseEncxWithClientID(t, c.ID)
		require.NoError(t, investigationHelpers.InsertCaseEncx(t, ctx, testPool, testCase))

		revokedToken, tokenHash, err := token.GenerateRawToken()
		require.NoError(t, err)
		revokedAt := time.Now().Add(-1 * time.Hour)
		insertTestToken(t, ctx, &token.Token{
			ID:        uuid.New(),
			CaseID:    testCase.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			RevokedAt: &revokedAt,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		})

		resp, err := http.Get(fmt.Sprintf("%s/token/%s/documents/estimate/pdf", testServerURL, revokedToken))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		_ = rawToken
	})
}

func toDocumentType(s string) documentDomain.DocumentType {
	switch s {
	case "estimate":
		return documentDomain.DocumentTypeEstimate
	case "mandate":
		return documentDomain.DocumentTypeMandate
	case "contract":
		return documentDomain.DocumentTypeContract
	case "invoice":
		return documentDomain.DocumentTypeInvoice
	default:
		return documentDomain.DocumentTypeOther
	}
}
