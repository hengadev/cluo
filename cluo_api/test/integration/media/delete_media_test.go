package media_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	caseHelpers "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	mediaHelpers "github.com/hengadev/cluo_api/test/helpers/media"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteMedia(t *testing.T) {
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
		t.Run("Administrator deletes media successfully", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// Create and insert properly encrypted media
			testMedia := mediaHelpers.CreateEncryptedTestMedia(t, ctx, crypto, testCase.ID)
			mediaHelpers.InsertMediaEncx(t, ctx, testPool, testMedia)

			// Verify media exists
			count := mediaHelpers.CountMediaByCaseID(t, ctx, testPool, testCase.ID)
			assert.Equal(t, 1, count)

			// Delete media
			resp := mediaHelpers.DeleteMediaRequest(t, ctx, testServerURL, adminToken, testMedia.ID.String())
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNoContent, resp.StatusCode)

			// Verify media was deleted
			count = mediaHelpers.CountMediaByCaseID(t, ctx, testPool, testCase.ID)
			assert.Equal(t, 0, count)

			t.Logf("Successfully deleted media with ID: %s", testMedia.ID)
		})

		t.Run("Client can delete media", func(t *testing.T) {
			clientToken := tu.SetupClientUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// Create and insert properly encrypted media
			testMedia := mediaHelpers.CreateEncryptedTestMedia(t, ctx, crypto, testCase.ID)
			mediaHelpers.InsertMediaEncx(t, ctx, testPool, testMedia)

			resp := mediaHelpers.DeleteMediaRequest(t, ctx, testServerURL, clientToken, testMedia.ID.String())
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Media not found", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			nonExistentID := uuid.New().String()

			resp := mediaHelpers.DeleteMediaRequest(t, ctx, testServerURL, adminToken, nonExistentID)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Invalid media ID format", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			resp := mediaHelpers.DeleteMediaRequest(t, ctx, testServerURL, adminToken, "invalid-uuid")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Unauthorized - no token", func(t *testing.T) {
			resp := mediaHelpers.DeleteMediaRequest(t, ctx, testServerURL, "", uuid.New().String())
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}
