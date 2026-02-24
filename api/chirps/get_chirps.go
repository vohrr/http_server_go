package chirps

import (
	"net/http"

	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/models"
)

func (h *ChirpHandler) GetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := h.Cfg.Queries.GetChirps(r.Context())
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
