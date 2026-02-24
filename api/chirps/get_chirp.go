package chirps

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/models"
)

func (h *ChirpHandler) GetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		api.RespondWithError(w, 401, "Bad Request")
		return
	}

	chirp, err := h.Cfg.Queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			api.RespondWithError(w, 404, "Not Found")
		} else {
			api.RespondWithError(w, 500, "Something went wrong")
		}
		return
	}
	api.RespondWithJSON(w, 200, models.MapChirpModel(chirp))
}
