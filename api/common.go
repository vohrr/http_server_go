package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vohrr/http_server_go/internal/auth"
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

func Authenticate(next http.HandlerFunc, cfg *ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//add authetication code here

		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}
		userID, err := auth.ValidateJWT(bearer, cfg.Secret)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}
		// if _, err := cfg.Queries.GetUserById(r.Context(), userID); err != nil {
		// 	RespondWithError(w, 401, "Not Authorized")
		// 	return
		// }
		ctx := context.WithValue(r.Context(), "userID", userID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
