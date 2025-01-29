package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	dbQueries 		*database.Queries
	platform		string
	tokenSecret 	string
}


func main() {
	godotenv.Load()

	// fetch database link url from environment
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Could not connect to database")
		return
	}

	// fetch platform type from environment
	currentPlatform := os.Getenv("PLATFORM")
	secretString := os.Getenv("TOKEN_SECRET_STRING")

	serveMux := http.NewServeMux()
	var server http.Server
	var cfg apiConfig

	cfg.dbQueries = database.New(db)
	cfg.platform = currentPlatform
	cfg.tokenSecret = secretString

	server.Addr = ":8080"
	server.Handler = serveMux

	serveMux.Handle(
		"/app/",
		cfg.middlewareMetricsInc(
			http.StripPrefix(
				"/app",http.FileServer(http.Dir(".")))))

	// API Routes
	serveMux.HandleFunc("GET /api/healthz", readinessCheck)

	// API User Routes
	serveMux.HandleFunc("POST /api/users", cfg.createUser)
	serveMux.HandleFunc("POST /api/login", cfg.loginUser)
	serveMux.HandleFunc("POST /api/refresh", cfg.refreshAccessToken)
	serveMux.HandleFunc("POST /api/revoke", cfg.revokeRefreshToken)

	// API Chirp Routes
	serveMux.HandleFunc("POST /api/chirps", cfg.createChirp)
	serveMux.HandleFunc("GET /api/chirps", cfg.getAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)

	// Admin Routes
	serveMux.HandleFunc("GET /admin/metrics", cfg.returnMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.deleteAllUsers)

	err = server.ListenAndServe()

	if (err != nil) {
		log.Fatal(err)
	}

}


