package admin

import (
	"fmt"
	"net/http"

	"github.com/vohrr/http_server_go/api"
)

type AdminHandler struct {
	Cfg *api.ApiConfig
}

func (h *AdminHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, h.Cfg.FileServerHits.Load())
}

func (h *AdminHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if h.Cfg.Platform != "dev" {
		w.WriteHeader(403)
	} else {
		h.Cfg.FileServerHits.Store(0)
		err := h.Cfg.Queries.AdminReset(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Something went wrong")
		}
	}
}

func (h *AdminHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "OK")
}

func (h *AdminHandler) RegisterSiteHit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
