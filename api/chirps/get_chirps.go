package chirps

import (
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/models"
	"github.com/vohrr/http_server_go/internal/database"
)

func (h *ChirpHandler) GetChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	sort_by := r.URL.Query().Get("sort")
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
	//sort chirps
	sortChirps(sort_by, &chirps)
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

func sortChirps(sort_by string, chirps *[]database.Chirp) {
	if len(sort_by) == 0 || sort_by == "asc" {
		slices.SortFunc(*chirps, func(a, b database.Chirp) int {
			if a.CreatedAt.Equal(b.CreatedAt) {
				return 0
			}
			if a.CreatedAt.After(b.CreatedAt) {
				return 1
			} else if b.CreatedAt.After(a.CreatedAt) {
				return -1
			}
			return 0
		})
	} else {
		slices.SortFunc(*chirps, func(a, b database.Chirp) int {
			if a.CreatedAt.Equal(b.CreatedAt) {
				return 0
			}
			if b.CreatedAt.After(a.CreatedAt) {
				return 1
			} else if a.CreatedAt.After(b.CreatedAt) {
				return -1
			}
			return 0
		})
	}
}
