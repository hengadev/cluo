package clientHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	"github.com/hengadev/cluo_api/internal/domain/client"
	clientHandler "github.com/hengadev/cluo_api/internal/interface/client"
	"github.com/stretchr/testify/require"
)

// NewCreateClientRequest creates an HTTP request for creating a client
func NewCreateClientRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	request client.CreateClientRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		serverURL+clientHandler.CreateClientEndpoint,
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

// NewGetClientByIDRequest creates an HTTP request for getting a client by ID
func NewGetClientByIDRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	clientID string,
	accessToken string,
) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		serverURL+clientHandler.ClientBasePath+"/"+clientID,
		nil,
	)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewDeleteClientRequest creates an HTTP request for deleting a client by ID
func NewDeleteClientRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	clientID string,
	accessToken string,
) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		serverURL+clientHandler.ClientBasePath+"/"+clientID,
		nil,
	)
	require.NoError(t, err)

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}
