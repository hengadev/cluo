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

// NewCreateContactRequest creates an HTTP request for creating a contact
func NewCreateContactRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	clientID string,
	request client.CreateContactRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		serverURL+clientHandler.ClientBasePath+"/"+clientID+clientHandler.ContactBasePath,
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

// NewGetContactByIDRequest creates an HTTP request for getting a contact by ID
func NewGetContactByIDRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	contactID string,
	accessToken string,
) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		serverURL+clientHandler.GetContactByIDEndpoint,
		nil,
	)
	require.NoError(t, err)

	// Replace the {id} placeholder with the actual contact ID
	req.URL.Path = clientHandler.ContactBasePath + "/" + contactID

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewGetContactIDsForClientRequest creates an HTTP request for getting contact IDs for a client
func NewGetContactIDsForClientRequest(
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
		serverURL+clientHandler.GetContactIDsForClientEndpoint,
		nil,
	)
	require.NoError(t, err)

	// Replace the {id} placeholder with the actual client ID
	req.URL.Path = clientHandler.ClientBasePath + "/" + clientID + "/contact-ids"

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewDeleteContactRequest creates an HTTP request for deleting a contact
func NewDeleteContactRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	contactID string,
	accessToken string,
) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		serverURL+clientHandler.DeleteContactEndpoint,
		nil,
	)
	require.NoError(t, err)

	// Replace the {id} placeholder with the actual contact ID
	req.URL.Path = clientHandler.ContactBasePath + "/" + contactID

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewUpdateContactRequest creates an HTTP request for updating a contact
func NewUpdateContactRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	contactID string,
	request client.UpdateContactRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		serverURL+clientHandler.UpdateContactEndpoint,
		bytes.NewReader(jsonBody),
	)
	require.NoError(t, err)

	// Replace the {id} placeholder with the actual contact ID
	req.URL.Path = clientHandler.ContactBasePath + "/" + contactID

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

// NewGetAllContactsByClientIDRequest creates an HTTP request for getting all contacts for a client
func NewGetAllContactsByClientIDRequest(
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
		serverURL+clientHandler.GetAllContactsByClientIDEndpoint,
		nil,
	)
	require.NoError(t, err)

	// Replace the {id} placeholder with the actual client ID
	req.URL.Path = clientHandler.ClientBasePath + "/" + clientID + clientHandler.ContactBasePath

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}
