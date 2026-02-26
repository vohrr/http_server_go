package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/auth"
	"github.com/vohrr/http_server_go/internal/database"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
	models.User
}

type LoginHandler struct {
	Cfg *ApiConfig
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials loginRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)
	if err != nil {
		RespondWithError(w, 400, "Bad Request")
		return
	}
	user, err := h.Cfg.Queries.GetUserByEmail(r.Context(), credentials.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, 401, "Incorrect email or password")

		} else {
			RespondWithError(w, 500, "Something went wrong")
		}
		return
	}
	if ok, _ := auth.CheckPasswordHash(credentials.Password, user.HashedPassword); !ok {
		RespondWithError(w, 401, "Incorrect email or password")
		return
	}

	expiresIn := time.Hour * 1
	jwt, err := auth.MakeJWT(user.ID, h.Cfg.Secret, expiresIn)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	refresh_token := auth.MakeRefreshToken()

	params := database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	}
	refresh, err := h.Cfg.Queries.CreateRefreshToken(r.Context(), params)
	if err != nil {
		RespondWithError(w, 500, "Error creating refresh token")
		return
	}
	response := loginResponse{
		User:         models.MapUserModel(user),
		Token:        jwt,
		RefreshToken: refresh.Token,
	}

	RespondWithJSON(w, 200, response)
}
