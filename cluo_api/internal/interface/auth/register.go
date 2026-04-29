package authHandler

import (
	"encoding/json"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/user"
)

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	cookies.SetAccessTokenCookie(w, resp.AccessToken, resp.AccessTokenExpiry)
	cookies.SetRefreshTokenCookie(w, resp.RefreshToken, resp.RefreshTokenExpiry)

	httpx.RespondWithJSON(w, resp, http.StatusCreated)
}
