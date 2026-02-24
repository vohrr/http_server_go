package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	RespondWithJSON(w, 200, models.MapUserModel(user))
}
