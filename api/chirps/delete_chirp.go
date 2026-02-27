package chirps

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
)

func (h *ChirpHandler) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	userIDString := r.Context().Value("userID").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		api.RespondWithJSON(w, 401, "Invalid userId")
		return
	}
	chirpID, _ := uuid.Parse(r.PathValue("chirpID"))
	chirp, err := h.Cfg.Queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			api.RespondWithJSON(w, 404, "Not found")
		} else {
			api.RespondWithJSON(w, 500, "Something went wrong")
		}
		return
	}
	if chirp.UserID != userID {
		api.RespondWithJSON(w, 403, "Forbidden")
		return
	}
	err = h.Cfg.Queries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		api.RespondWithJSON(w, 500, "Something went wrong")
		return
	}

	w.WriteHeader(204)
}
