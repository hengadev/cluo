package authHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *AuthHandler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	resp, err := h.svc.RefreshSession(ctx, sessionInfo.ID)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	cookies.SetAccessTokenCookie(w, resp.AccessToken, resp.AccessTokenExpiry)
	cookies.SetRefreshTokenCookie(w, resp.RefreshToken, resp.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, resp, http.StatusOK)
}
