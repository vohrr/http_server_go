package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/internal/database"
)

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	//decode the request body
	type createUserRequest struct {
		Email string `json:"email"`
	}
	var data createUserRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}
	var user database.User
	//write the request to the db
	user, err = h.Cfg.Queries.CreateUser(r.Context(), data.Email)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}
	//write user to responsewriter
	api.RespondWithJSON(w, 201, mapUserModel(user))
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func mapUserModel(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}
