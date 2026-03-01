package chirps

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/database"
)

func (h *ChirpHandler) GetChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	if len(author) != 0 {
		userID, err := uuid.Parse(author)
		if err != nil {
			api.RespondWithError(w, 400, "Bad Request")
			return
		}
		chirps, err = h.Cfg.Queries.GetChirpsByUser(r.Context(), userID)
	} else {
		chirps, err = h.Cfg.Queries.GetChirps(r.Context())
	}

	if err != nil {
		api.RespondWithError(w, 500, "Could not fetch chirp data.")
		return
	}
	var chirpsResponse []models.Chirp
	for _, chirp := range chirps {
		chirpsResponse = append(chirpsResponse, models.MapChirpModel(chirp))
	}
	api.RespondWithJSON(w, 200, chirpsResponse)
}
