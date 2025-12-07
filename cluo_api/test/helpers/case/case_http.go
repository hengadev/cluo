package caseHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	"github.com/stretchr/testify/require"
)

// NewCreateCaseRequest creates an HTTP request for creating a case
func NewCreateCaseRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	request caseDomain.CreateCaseRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		serverURL+"/cases",
		bytes.NewReader(jsonBody),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

