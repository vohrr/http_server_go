package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vohrr/http_server_go/api"
	"github.com/vohrr/http_server_go/api/admin"
	"github.com/vohrr/http_server_go/api/chirps"
	"github.com/vohrr/http_server_go/api/users"
	"github.com/vohrr/http_server_go/internal/database"
)

func main() {
	//init db connection
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	//server config
	mux := http.NewServeMux()
	applicationPath := "/app/"
	appHandler := http.StripPrefix(applicationPath, http.FileServer(http.Dir(".")))
	//register handlers
	userHandler := users.UserHandler{Cfg: cfg}
	adminHandler := admin.AdminHandler{Cfg: cfg}
	chirpsHandler := chirps.ChirpHandler{Cfg: cfg}
	loginHandler := api.LoginHandler{Cfg: cfg}
	//register endpoints
	mux.Handle(applicationPath, adminHandler.RegisterSiteHit(appHandler))
	mux.HandleFunc("GET /admin/metrics", adminHandler.Metrics)
	mux.HandleFunc("GET /api/healthz", adminHandler.Health)
	mux.HandleFunc("POST /admin/reset", adminHandler.Reset)
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("PUT /api/users", api.Authenticate(userHandler.UpdateUser, cfg))
	mux.HandleFunc("POST /api/chirps", api.Authenticate(chirpsHandler.CreateChirp, cfg))
	mux.HandleFunc("GET /api/chirps", chirpsHandler.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirpsHandler.GetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", api.Authenticate(chirpsHandler.DeleteChirp, cfg))
	mux.HandleFunc("POST /api/login", loginHandler.Login)
	mux.HandleFunc("POST /api/refresh", loginHandler.Refresh)
	mux.HandleFunc("POST /api/revoke", loginHandler.Revoke)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080", //inject configurable host address
	}

	log.Fatal(server.ListenAndServe())
}

func loadConfig() (*api.ApiConfig, error) {
	var cfg api.ApiConfig
	err := godotenv.Load()
	if err != nil {
		return &cfg, err
	}
	cfg.Platform = os.Getenv("PLATFORM")
	cfg.Queries, err = initDatabase()
	cfg.Secret = os.Getenv("SECRET")
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func initDatabase() (*database.Queries, error) {
	dbConnectionString := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("Unable to establish connection to database")
	}
	return database.New(db), nil
}
