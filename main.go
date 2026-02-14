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
	mux.HandleFunc("/metrics", cfg.Metrics)
	mux.HandleFunc("/reset", cfg.Reset)
	mux.HandleFunc("/healthz", cfg.Health)
	//mux.Handle("/healthz", HealthzHandler{})

	server := http.Server{
		Handler: mux,
		Addr:    ":8080", //inject configurable host address
	}

	log.Fatal(server.ListenAndServe())
}
