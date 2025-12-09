package media_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	caseHelpers "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	mediaHelpers "github.com/hengadev/cluo_api/test/helpers/media"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	mediaDomain "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMediaByCaseID(t *testing.T) {
	ctx := context.Background()

	// Helper function to setup a client
	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		c := clientHelpers.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return c.ID
	}

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("List all media for a case", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// Create and insert multiple media files with proper encryption
			for i := 0; i < 5; i++ {
				media := mediaHelpers.CreateEncryptedTestMedia(t, ctx, crypto, testCase.ID)
				mediaHelpers.InsertMediaEncx(t, ctx, testPool, media)
			}

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, adminToken, testCase.ID.String(), 1, 20, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.ListMediaResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Len(t, response.Media, 5)
			assert.Equal(t, 5, response.Pagination.TotalItems)
			assert.Equal(t, 1, response.Pagination.Page)
			assert.Equal(t, 20, response.Pagination.PageSize)

			t.Logf("Successfully listed %d media files for case: %s", len(response.Media), testCase.ID.String())
		})

		t.Run("Pagination - first page", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// Create and insert 15 media files with proper encryption
			for i := 0; i < 15; i++ {
				media := mediaHelpers.CreateEncryptedTestMedia(t, ctx, crypto, testCase.ID)
				mediaHelpers.InsertMediaEncx(t, ctx, testPool, media)
			}

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, adminToken, testCase.ID.String(), 1, 10, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.ListMediaResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Len(t, response.Media, 10)
			assert.Equal(t, 15, response.Pagination.TotalItems)
			assert.Equal(t, 2, response.Pagination.TotalPages)
		})

		t.Run("Pagination - second page", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// Create and insert 15 media files with proper encryption
			for i := 0; i < 15; i++ {
				media := mediaHelpers.CreateEncryptedTestMedia(t, ctx, crypto, testCase.ID)
				mediaHelpers.InsertMediaEncx(t, ctx, testPool, media)
			}

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, adminToken, testCase.ID.String(), 2, 10, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.ListMediaResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Len(t, response.Media, 5) // Remaining items
			assert.Equal(t, 15, response.Pagination.TotalItems)
		})

		t.Run("Empty results - no media for case", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, adminToken, testCase.ID.String(), 1, 20, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.ListMediaResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Len(t, response.Media, 0)
			assert.Equal(t, 0, response.Pagination.TotalItems)
		})

		t.Run("Client can list media", func(t *testing.T) {
			clientToken := tu.SetupClientUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// Create and insert properly encrypted media
			media := mediaHelpers.CreateEncryptedTestMedia(t, ctx, crypto, testCase.ID)
			mediaHelpers.InsertMediaEncx(t, ctx, testPool, media)

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, clientToken, testCase.ID.String(), 1, 20, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("Invalid case ID format", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, adminToken, "invalid-uuid", 1, 20, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Invalid page number", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, adminToken, testCase.ID.String(), 0, 20, nil)
			defer resp.Body.Close()

			// Should default to page 1 or return error
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)
		})
	})

	t.Run("Authentication", func(t *testing.T) {
		t.Run("No authentication token", func(t *testing.T) {
			resp := mediaHelpers.ListMediaByCaseIDRequest(t, ctx, testServerURL, "", uuid.New().String(), 1, 20, nil)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}
