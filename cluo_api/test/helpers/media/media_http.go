package mediaHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/stretchr/testify/require"
)

// UploadMediaRequest sends POST /media with multipart form data
func UploadMediaRequest(
	t *testing.T,
	ctx context.Context,
	baseURL string,
	token string,
	fileContent []byte,
	fileName string,
	mimeType string,
	caseID string,
	caption string,
	isPublished string,
) *http.Response {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file with proper MIME type
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName)}
	h["Content-Type"] = []string{mimeType}
	part, err := writer.CreatePart(h)
	require.NoError(t, err)
	_, err = io.Copy(part, bytes.NewReader(fileContent))
	require.NoError(t, err)

	// Add form fields
	_ = writer.WriteField("caseId", caseID)
	if caption != "" {
		_ = writer.WriteField("caption", caption)
	}
	if isPublished != "" {
		_ = writer.WriteField("isPublished", isPublished)
	}

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		baseURL+"/media",
		body,
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	if token != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: token,
		}
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// GetMediaByIDRequest sends GET /media/{id}
func GetMediaByIDRequest(t *testing.T, ctx context.Context, baseURL string, token string, mediaID string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/media/%s", baseURL, mediaID),
		nil,
	)
	require.NoError(t, err)

	if token != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: token,
		}
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// UpdateMediaRequest sends PATCH /media/{id}
func UpdateMediaRequest(
	t *testing.T,
	ctx context.Context,
	baseURL string,
	token string,
	mediaID string,
	payload *domain.UpdateMediaRequest,
) *http.Response {
	t.Helper()

	bodyBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		"PATCH",
		fmt.Sprintf("%s/media/%s", baseURL, mediaID),
		bytes.NewReader(bodyBytes),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: token,
		}
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// DeleteMediaRequest sends DELETE /media/{id}
func DeleteMediaRequest(t *testing.T, ctx context.Context, baseURL string, token string, mediaID string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/media/%s", baseURL, mediaID),
		nil,
	)
	require.NoError(t, err)

	if token != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: token,
		}
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// ListMediaByCaseIDRequest sends GET /case/{caseId}/media
func ListMediaByCaseIDRequest(
	t *testing.T,
	ctx context.Context,
	baseURL string,
	token string,
	caseID string,
	page, pageSize int,
	mediaType *string,
) *http.Response {
	t.Helper()

	url := fmt.Sprintf("%s/case/%s/media?page=%d&pageSize=%d", baseURL, caseID, page, pageSize)
	if mediaType != nil {
		url += "&type=" + *mediaType
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	if token != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: token,
		}
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
