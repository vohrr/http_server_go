package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	//register request handlers
	mux.Handle(applicationPath, cfg.RegisterSiteHit(appHandler))
	mux.HandleFunc("GET /admin/metrics", cfg.Metrics)
	mux.HandleFunc("GET /api/healthz", cfg.Health)
	mux.HandleFunc("POST /admin/reset", cfg.Reset)
	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirp)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080", //inject configurable host address
	}

	log.Fatal(server.ListenAndServe())
}

func loadConfig() (*apiConfig, error) {
	var cfg apiConfig
	err := godotenv.Load()
	if err != nil {
		return &cfg, err
	}
	cfg.platform = os.Getenv("PLATFORM")
	cfg.queries, err = initDatabase()
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
