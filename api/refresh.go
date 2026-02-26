package api

import (
	"net/http"
	"time"

	"github.com/vohrr/http_server_go/internal/auth"
)

func (h *LoginHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Invalid Bearer token")
		return
	}
	refresh_token, err := h.Cfg.Queries.GetRefreshToken(r.Context(), bearer)
	if err != nil {
		RespondWithError(w, 401, "Invalid Bearer token")
		return
	}

	if refresh_token.RevokedAt.Valid {
		RespondWithError(w, 401, "Invalid Bearer token")
		return
	}

	newJWT, err := auth.MakeJWT(refresh_token.UserID, h.Cfg.Secret, time.Hour)
	if err != nil {
		RespondWithError(w, 500, "Error generating JWT")
		return
	}

	payload := struct {
		Token string `json:"token"`
	}{
		Token: newJWT,
	}
	RespondWithJSON(w, 200, payload)
}
