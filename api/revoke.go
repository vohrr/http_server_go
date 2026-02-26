package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/vohrr/http_server_go/internal/auth"
	"github.com/vohrr/http_server_go/internal/database"
)

func (h *LoginHandler) Revoke(w http.ResponseWriter, r *http.Request) {
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
	revokedTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	params := database.UpdateRefreshTokenParams{
		RevokedAt: revokedTime,
		Token:     refresh_token.Token,
	}
	err = h.Cfg.Queries.UpdateRefreshToken(r.Context(), params)
	if err != nil {
		RespondWithError(w, 500, "Unable to revoke refresh token")
		return
	}

	w.WriteHeader(204)
}
