package authHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
)

func (h *handler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	err := h.svc.SignOut(ctx, sessionInfo)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	cookies.ClearTokenCookies(w)

	httpx.RespondWithJSON(w, map[string]string{"message": "signed out successfully"}, http.StatusOK)
}
