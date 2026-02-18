package main

import (
	"encoding/json"
	"fmt"
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

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	data, _ := json.Marshal(payload)
	w.WriteHeader(statusCode)
	w.Write(data)
}
