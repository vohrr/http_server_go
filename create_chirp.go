package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/internal/database"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	var c createChirpRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if status, err := validateChirp(&c); err != nil {
		respondWithError(w, status, err.Error())
		return
	}

	params := database.CreateChirpParams{
		Body:   c.Body,
		UserID: c.UserId,
	}
	chirp, err := cfg.queries.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	respondWithJSON(w, 201, mapChirpModel(chirp))
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

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func mapChirpModel(db database.Chirp) Chirp {
	return Chirp{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		Body:      db.Body,
		UserId:    db.UserID,
	}
}
