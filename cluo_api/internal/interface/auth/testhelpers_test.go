package authHandler

import (
	"context"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignIn(ctx context.Context, req *user.SignInRequest) (*user.CreateSessionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.CreateSessionResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, req *user.RegisterRequest) (*user.CreateSessionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.CreateSessionResponse), args.Error(1)
}

func (m *MockAuthService) SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error {
	args := m.Called(ctx, sessionInfo)
	return args.Error(0)
}

func (m *MockAuthService) RefreshSession(ctx context.Context, sessionID uuid.UUID) (*user.RefreshSessionResponse, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.RefreshSessionResponse), args.Error(1)
}

func (m *MockAuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*user.CurrentUserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.CurrentUserResponse), args.Error(1)
}

type MockAuthMiddleware struct {
	mock.Mock
}

func (m *MockAuthMiddleware) RequireMinimumRole(role identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}
}

func (m *MockAuthMiddleware) RequireAnyRole(roles ...identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}
}

func (m *MockAuthMiddleware) RequireAccessToken(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func (m *MockAuthMiddleware) RequireRefreshToken(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func (m *MockAuthMiddleware) RequireAdmin(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func (m *MockAuthMiddleware) RequireServiceAuth(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
