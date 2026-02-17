package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	var cfg apiConfig
	applicationPath := "/app/"
	appHandler := http.StripPrefix(applicationPath, http.FileServer(http.Dir(".")))
	//register request handlers
	mux.Handle(applicationPath, cfg.RegisterSiteHit(appHandler))
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.Metrics)
	mux.HandleFunc("GET /api/healthz", cfg.Health)
	mux.HandleFunc("POST /admin/reset", cfg.Reset)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080", //inject configurable host address
	}

	log.Fatal(server.ListenAndServe())
}
