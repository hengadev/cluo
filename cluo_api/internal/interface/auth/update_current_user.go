package authHandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/user"
)

func (h *AuthHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionInfo, ok := session.SessionInfoFromContext(ctx)
	if !ok {
		httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req user.UpdateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		httpx.RespondWithError(w, errors.New("name is required"), http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateCurrentUserName(ctx, sessionInfo.UserID, req.Name); err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
