package main

import (
	"net/http"
)

func (cfg *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.queries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Could not fetch chirp data.")
		return
	}
	var chirpsResponse []Chirp
	for _, chirp := range chirps {
		chirpsResponse = append(chirpsResponse, mapChirpModel(chirp))
	}
	respondWithJSON(w, 200, chirpsResponse)
}
