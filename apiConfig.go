package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	data := fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())
	body, err := json.Marshal(&data)
	if r.Method != "GET" || err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write(body)
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
}

func (cfg *apiConfig) RegisterSiteHit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) Health(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal("OK")
	if r.Method != "GET" || err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write(body)
}
