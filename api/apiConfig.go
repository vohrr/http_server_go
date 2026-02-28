package api

import (
	"sync/atomic"

	"github.com/vohrr/http_server_go/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
	Platform       string
	Secret         string
	PolkaAPIKey    string
}
