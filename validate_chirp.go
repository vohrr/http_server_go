package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type chirp struct {
	Body string `json:"body"`
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	var c chirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(c.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	payload := struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	}
	respondWithJSON(w, 200, payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	data, _ := json.Marshal(payload)
	w.WriteHeader(statusCode)
	w.Write(data)
}
