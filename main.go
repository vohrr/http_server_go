package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	var cfg apiConfig
	applicationPath := "/app/"
	basicHandler := http.StripPrefix(applicationPath, http.FileServer(http.Dir(".")))

	//register request handlers
	mux.Handle(applicationPath, cfg.RegisterSiteHit(basicHandler))
	mux.HandleFunc("GET /api/metrics", cfg.Metrics)
	mux.HandleFunc("GET /api/healthz", cfg.Health)
	mux.HandleFunc("POST /api/reset", cfg.Reset)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080", //inject configurable host address
	}

	log.Fatal(server.ListenAndServe())
}
