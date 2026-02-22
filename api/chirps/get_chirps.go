package chirps

import (
	"net/http"

	"github.com/vohrr/http_server_go/api"
)

func (h *ChirpHandler) GetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := h.Cfg.Queries.GetChirps(r.Context())
	if err != nil {
		api.RespondWithError(w, 500, "Could not fetch chirp data.")
		return
	}
	var chirpsResponse []Chirp
	for _, chirp := range chirps {
		chirpsResponse = append(chirpsResponse, mapChirpModel(chirp))
	}
	api.RespondWithJSON(w, 200, chirpsResponse)
}
