package media_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	caseHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	mediaHelpers "github.com/hengadev/cluo_api/test/helpers/media"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	mediaDomain "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMedia(t *testing.T) {
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
		t.Run("Update caption only", func(t *testing.T) {
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

			newCaption := "Updated caption"
			updateRequest := &mediaDomain.UpdateMediaRequest{
				Caption: &newCaption,
			}

			resp := mediaHelpers.UpdateMediaRequest(t, ctx, testServerURL, adminToken, testMedia.ID.String(), updateRequest)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.MediaResponse
			updateErr := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, updateErr)

			assert.Equal(t, newCaption, response.Caption)
			assert.Equal(t, testMedia.ID.String(), response.ID)

			t.Logf("Successfully updated caption for media: %s", response.ID)
		})

		t.Run("Update isPublished only", func(t *testing.T) {
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

			isPublished := true
			updateRequest := &mediaDomain.UpdateMediaRequest{
				IsPublished: &isPublished,
			}

			resp := mediaHelpers.UpdateMediaRequest(t, ctx, testServerURL, adminToken, testMedia.ID.String(), updateRequest)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.MediaResponse
			updateErr := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, updateErr)

			assert.True(t, response.IsPublished)
		})

		t.Run("Update both caption and isPublished", func(t *testing.T) {
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

			newCaption := "New caption"
			isPublished := true
			updateRequest := &mediaDomain.UpdateMediaRequest{
				Caption:     &newCaption,
				IsPublished: &isPublished,
			}

			resp := mediaHelpers.UpdateMediaRequest(t, ctx, testServerURL, adminToken, testMedia.ID.String(), updateRequest)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response mediaDomain.MediaResponse
			updateErr := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, updateErr)

			assert.Equal(t, newCaption, response.Caption)
			assert.True(t, response.IsPublished)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("Caption too long", func(t *testing.T) {
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

			longCaption := string(make([]byte, 501)) // > 500 chars
			updateRequest := &mediaDomain.UpdateMediaRequest{
				Caption: &longCaption,
			}

			resp := mediaHelpers.UpdateMediaRequest(t, ctx, testServerURL, adminToken, testMedia.ID.String(), updateRequest)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Media not found", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			newCaption := "New caption"
			updateRequest := &mediaDomain.UpdateMediaRequest{
				Caption: &newCaption,
			}

			resp := mediaHelpers.UpdateMediaRequest(t, ctx, testServerURL, adminToken, uuid.New().String(), updateRequest)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Authentication", func(t *testing.T) {
		t.Run("No authentication token", func(t *testing.T) {
			newCaption := "New caption"
			updateRequest := &mediaDomain.UpdateMediaRequest{
				Caption: &newCaption,
			}

			resp := mediaHelpers.UpdateMediaRequest(t, ctx, testServerURL, "", uuid.New().String(), updateRequest)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}
