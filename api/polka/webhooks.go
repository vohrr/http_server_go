package polka

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/internal/database"
)

type WebhookHandler struct {
	Cfg *api.ApiConfig
}

type WebhookEventRequest struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data"`
}

func (h *WebhookHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	var req WebhookEventRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&req)

	if err != nil {
		api.RespondWithError(w, 400, "Bad Request")
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	if userID, ok := req.Data["user_id"]; !ok {
		api.RespondWithError(w, 400, "Bad Request")
		return
	} else {
		uId, _ := uuid.Parse(userID.(string))
		params := database.SetChirpyRedParams{
			IsChirpyRed: true,
			ID:          uId,
		}
		err = h.Cfg.Queries.SetChirpyRed(r.Context(), params)
		if err != nil {
			api.RespondWithError(w, 404, "Someting went wrong")
			return
		}

		w.WriteHeader(204)
	}
}
