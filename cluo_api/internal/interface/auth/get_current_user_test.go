package authHandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentUser_Unauthorized(t *testing.T) {
	// Setup
	mockSvc := new(MockAuthService)
	mockAuthMw := new(MockAuthMiddleware)
	h := New(mockSvc, mockAuthMw)

	// Create request without session context
	req := httptest.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()

	// Execute
	h.(*AuthHandler).GetCurrentUser(w, req)

	// Assert
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "unauthorized")
}

func TestGetCurrentUser_Success(t *testing.T) {
	// Setup
	userID := uuid.New()
	expectedEmail := "test@example.com"
	expectedRole := "client"

	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: userID,
		Role:   identity.Client,
		State:  session.SessionActive,
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), sessionInfo)

	expectedResponse := &user.CurrentUserResponse{
		ID:    userID.String(),
		Email: expectedEmail,
		Role:  expectedRole,
	}

	mockSvc := new(MockAuthService)
	mockSvc.On("GetCurrentUser", ctx, userID).Return(expectedResponse, nil)

	mockAuthMw := new(MockAuthMiddleware)
	h := New(mockSvc, mockAuthMw)

	// Create request with session context
	req := httptest.NewRequest("GET", "/auth/me", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	h.(*AuthHandler).GetCurrentUser(w, req)

	// Assert
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response user.CurrentUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, userID.String(), response.ID)
	assert.Equal(t, expectedEmail, response.Email)
	assert.Equal(t, expectedRole, response.Role)

	mockSvc.AssertExpectations(t)
}

func TestGetCurrentUser_ServiceError(t *testing.T) {
	// Setup
	userID := uuid.New()
	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: userID,
		Role:   identity.Client,
		State:  session.SessionActive,
	}

	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), sessionInfo)

	mockSvc := new(MockAuthService)
	mockSvc.On("GetCurrentUser", ctx, userID).Return(nil, assert.AnError)

	mockAuthMw := new(MockAuthMiddleware)
	h := New(mockSvc, mockAuthMw)

	// Create request with session context
	req := httptest.NewRequest("GET", "/auth/me", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	h.(*AuthHandler).GetCurrentUser(w, req)

	// Assert
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockSvc.AssertExpectations(t)
}
