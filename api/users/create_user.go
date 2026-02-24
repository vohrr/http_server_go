package users

import (
	"encoding/json"
	"net/http"

	"github.com/vohrr/http_server_go/api"

	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/auth"
	"github.com/vohrr/http_server_go/internal/database"
)

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	//decode the request body
	var data createUserRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}
	if !data.isValid() {
		api.RespondWithError(w, 400, "Bad Request")
	}
	var user database.User
	hashedPassword, err := auth.HashPassword(data.Password)
	userParams := database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          data.Email,
	}
	//write the request to the db
	user, err = h.Cfg.Queries.CreateUser(r.Context(), userParams)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}
	//write user to responsewriter
	api.RespondWithJSON(w, 201, models.MapUserModel(user))
}

type createUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r createUserRequest) isValid() bool {
	return (len(r.Email) > 0 && len(r.Password) > 0)
}
