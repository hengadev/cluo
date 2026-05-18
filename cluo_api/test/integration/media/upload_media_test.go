package media_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

func TestUploadMedia(t *testing.T) {
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
		t.Run("Administrator uploads image successfully", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			// Create a test client first
			clientID := setupClient(t, ctx)

			// Create a test case
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			insertErr := caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)
			require.NoError(t, insertErr, "Failed to insert test case")

			// Create test image data
			imageData := []byte("fake-image-content-for-testing")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				imageData,
				"test-image.jpg",
				"image/jpeg",
				testCase.ID.String(),
				"Test caption",
				"false",
			)
			defer resp.Body.Close()

			// Read and log response body
			bodyBytes, readErr := io.ReadAll(resp.Body)
			require.NoError(t, readErr)
			t.Logf("Response status: %d, body: %s", resp.StatusCode, string(bodyBytes))

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var response mediaDomain.MediaResponse
			decodeErr := json.Unmarshal(bodyBytes, &response)
			require.NoError(t, decodeErr, "Failed to decode response")

			assert.NotEmpty(t, response.ID)
			assert.Equal(t, testCase.ID.String(), response.CaseID)
			assert.Equal(t, "test-image.jpg", response.FileName)
			assert.Equal(t, "image/jpeg", response.MimeType)
			assert.Equal(t, "image", response.Type)
			assert.Equal(t, int64(len(imageData)), response.FileSize)
			assert.Equal(t, "Test caption", response.Caption)
			assert.False(t, response.IsPublished)
			assert.NotEmpty(t, response.URL)

			t.Logf("Successfully uploaded media with ID: %s", response.ID)
		})

		t.Run("Administrator uploads video successfully", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			videoData := []byte("fake-video-content-for-testing")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				videoData,
				"test-video.mp4",
				"video/mp4",
				testCase.ID.String(),
				"",
				"true",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var response mediaDomain.MediaResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, "video", response.Type)
			assert.Equal(t, "video/mp4", response.MimeType)
			assert.True(t, response.IsPublished)
		})

		t.Run("Administrator uploads audio successfully", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			audioData := []byte("fake-audio-content-for-testing")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				audioData,
				"test-audio.mp3",
				"audio/mpeg",
				testCase.ID.String(),
				"Audio recording",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var response mediaDomain.MediaResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, "audio", response.Type)
			assert.Equal(t, "audio/mpeg", response.MimeType)
			assert.Equal(t, "Audio recording", response.Caption)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("Missing case ID", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			imageData := []byte("test-content")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				imageData,
				"test.jpg",
				"image/jpeg",
				"", // Empty case ID
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Invalid case ID format", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			imageData := []byte("test-content")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				imageData,
				"test.jpg",
				"image/jpeg",
				"invalid-uuid",
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Case not found", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)

			imageData := []byte("test-content")
			nonExistentCaseID := uuid.New().String()

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				imageData,
				"test.jpg",
				"image/jpeg",
				nonExistentCaseID,
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Unsupported file type", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			fileData := []byte("test-content")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				fileData,
				"test.pdf",
				"application/pdf",
				testCase.ID.String(),
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Authentication and Authorization", func(t *testing.T) {
		t.Run("No authentication token", func(t *testing.T) {
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			imageData := []byte("test-content")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, "", // No token
				imageData,
				"test.jpg",
				"image/jpeg",
				testCase.ID.String(),
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Client role can upload media", func(t *testing.T) {
			clientToken := tu.SetupClientUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			imageData := []byte("test-content")

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, clientToken,
				imageData,
				"test.jpg",
				"image/jpeg",
				testCase.ID.String(),
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	})

	t.Run("File Size Limits", func(t *testing.T) {
		t.Run("Large file within limit", func(t *testing.T) {
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer caseHelpers.ClearCasesTable(t, ctx, testPool)
			defer clientHelpers.ClearClientsTable(t, ctx, testPool)
			defer mediaHelpers.ClearMediaTable(t, ctx, testPool)

			clientID := setupClient(t, ctx)
			testCase := caseHelpers.NewTestCaseEncxWithClientID(t, clientID)
			caseHelpers.InsertCaseEncx(t, ctx, testPool, testCase)

			// 5 MB file
			largeData := bytes.Repeat([]byte("x"), 5*1024*1024)

			resp := mediaHelpers.UploadMediaRequest(
				t, ctx, testServerURL, adminToken,
				largeData,
				"large-file.jpg",
				"image/jpeg",
				testCase.ID.String(),
				"",
				"",
			)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	})
}
