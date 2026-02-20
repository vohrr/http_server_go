package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/vohrr/http_server_go/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	queries        *database.Queries
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileServerHits.Load())
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
}

func (cfg *apiConfig) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "OK")
}

func (cfg *apiConfig) RegisterSiteHit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
