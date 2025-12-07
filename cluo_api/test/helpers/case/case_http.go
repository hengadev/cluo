package caseHelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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

