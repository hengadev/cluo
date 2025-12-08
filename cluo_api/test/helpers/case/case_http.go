package caseHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
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

// NewGetCaseByIDRequest creates an HTTP request for retrieving a case by ID
func NewGetCaseByIDRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	caseID uuid.UUID,
	accessToken string,
) *http.Request {
	t.Helper()

	url := serverURL + "/cases/" + caseID.String()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil, // No body for GET request
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

// NewDeleteCaseRequest creates an HTTP request for deleting a case by ID
func NewDeleteCaseRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	caseID uuid.UUID,
	accessToken string,
) *http.Request {
	t.Helper()

	url := serverURL + "/cases/" + caseID.String()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil, // No body for DELETE request
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

// NewUpdateCaseRequest creates an HTTP request for updating a case by ID
func NewUpdateCaseRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	caseID uuid.UUID,
	updateRequest *caseDomain.UpdateCaseRequest,
	accessToken string,
) *http.Request {
	t.Helper()

	url := serverURL + "/cases/" + caseID.String()

	// Marshal update request to JSON
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err, "Failed to marshal update request")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		bytes.NewReader(requestBody),
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

// NewListCasesRequest creates an HTTP request for listing cases with filtering and pagination
func NewListCasesRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	page int,
	pageSize int,
	filters map[string]string,
	accessToken string,
) *http.Request {
	t.Helper()

	// Build URL with query parameters
	url := serverURL + "/cases"

	// Create request with no body (GET request)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil, // No body for GET request
	)
	require.NoError(t, err)

	// Add query parameters
	q := req.URL.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))

	// Add filter parameters if provided
	for key, value := range filters {
		if value != "" {
			q.Set(key, value)
		}
	}
	req.URL.RawQuery = q.Encode()

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewListCasesByClientRequest creates an HTTP request for listing cases for a specific client
func NewListCasesByClientRequest(
	t *testing.T,
	ctx context.Context,
	serverURL string,
	clientID uuid.UUID,
	page int,
	pageSize int,
	accessToken string,
) *http.Request {
	t.Helper()

	url := serverURL + "/clients/" + clientID.String() + "/cases"

	// Create request with no body (GET request)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil, // No body for GET request
	)
	require.NoError(t, err)

	// Add pagination query parameters
	q := req.URL.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	req.URL.RawQuery = q.Encode()

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

