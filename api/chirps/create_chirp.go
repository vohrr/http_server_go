package chirps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/database"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (h *ChirpHandler) CreateChirp(w http.ResponseWriter, r *http.Request) {
	var c createChirpRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}

	if status, err := validateChirp(&c); err != nil {
		api.RespondWithError(w, status, err.Error())
		return
	}

	params := database.CreateChirpParams{
		Body:   c.Body,
		UserID: c.UserId,
	}
	chirp, err := h.Cfg.Queries.CreateChirp(r.Context(), params)
	if err != nil {
		api.RespondWithError(w, 500, "Something went wrong")
		return
	}

	api.RespondWithJSON(w, 201, models.MapChirpModel(chirp))
}

func validateChirp(c *createChirpRequest) (int, error) {
	//validate length
	if len(c.Body) > 140 {
		return 400, fmt.Errorf("Chirp is too long")
	}

	//filter profanity
	profanity := getProfanityMap()
	chirpSlice := strings.Split(c.Body, " ")
	for i, word := range chirpSlice {
		if _, ok := profanity[strings.ToLower(word)]; ok {
			chirpSlice[i] = "****"
		}
	}
	c.Body = strings.Join(chirpSlice, " ")
	return 201, nil
}

func getProfanityMap() map[string]any {
	return map[string]any{
		"kerfuffle": "",
		"sharbert":  "",
		"fornax":    "",
	}

}
