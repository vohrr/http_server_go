package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", message)
}
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	data, _ := json.Marshal(payload)
	w.WriteHeader(statusCode)
	w.Write(data)
}
