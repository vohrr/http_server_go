package users

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/auth"
	"github.com/vohrr/http_server_go/internal/database"
)

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var data UserRequest

	userId, ok := r.Context().Value("userID").(string)
	if !ok {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}
	userID, _ := uuid.Parse(userId)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		api.RespondWithError(w, 400, "Bad Request")
		return
	}
	hashed_pw, err := auth.HashPassword(data.Password)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}

	params := database.UpdateUserParams{
		HashedPassword: hashed_pw,
		Email:          data.Email,
		ID:             userID,
	}
	updatedUser, err := h.Cfg.Queries.UpdateUser(r.Context(), params)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}

	api.RespondWithJSON(w, 200, models.MapUserModel(updatedUser))

}
