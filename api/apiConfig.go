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
	platform       string
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
	if cfg.platform != "dev" {
		w.WriteHeader(403)
	} else {
		cfg.fileServerHits.Store(0)
		err := cfg.queries.AdminReset(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Something went wrong")
		}
	}
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
