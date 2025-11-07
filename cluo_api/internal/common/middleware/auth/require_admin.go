package auth

import (
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

// RequireAdmin validates session and ensures user has admin role
func (m *SessionAuthMiddleware) RequireAdmin(next mw.Handler) mw.Handler {
	return m.RequireMinimumRole(identity.Administrator)(next)
}
