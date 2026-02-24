package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vohrr/http_server_go/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func MapChirpModel(db database.Chirp) Chirp {
	return Chirp{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		Body:      db.Body,
		UserId:    db.UserID,
	}
}
