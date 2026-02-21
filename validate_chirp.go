package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

func getProfanityMap() map[string]any {
	return map[string]any{
		"kerfuffle": "",
		"sharbert":  "",
		"fornax":    "",
	}

}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	var c chirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	//validate length
	if len(c.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	//filter profanity
	profanity := getProfanityMap()
	chirpSlice := strings.Split(c.Body, " ")
	for i, word := range chirpSlice {
		if _, ok := profanity[strings.ToLower(word)]; ok {
			chirpSlice[i] = "****"
		}
	}

	payload := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: strings.Join(chirpSlice, " "),
	}
	respondWithJSON(w, 200, payload)
}
