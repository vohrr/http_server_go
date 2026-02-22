package main

import (
	"database/sql"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 401, "Bad Request")
		return
	}

	chirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Not Found")
		} else {
			respondWithError(w, 500, "Something went wrong")
		}
		return
	}
	respondWithJSON(w, 200, mapChirpModel(chirp))
}
