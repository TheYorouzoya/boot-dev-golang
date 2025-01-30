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
	polkaKey		string
}


func main() {
	godotenv.Load()

	// fetch database link url from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Could not get DB URL from the environment")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Could not connect to database")
		return
	}

	// fetch platform type from environment
	currentPlatform := os.Getenv("PLATFORM")
	if currentPlatform == "" {
		log.Fatal("Could not get playform variable from the environment")
	}

	secretString := os.Getenv("TOKEN_SECRET_STRING")
	if secretString == "" {
		log.Fatal("Could not get token generation secret string from the environment")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("Could not get polka key from the environment")
	}

	serveMux := http.NewServeMux()
	var server http.Server
	var cfg apiConfig

	cfg.dbQueries = database.New(db)
	cfg.platform = currentPlatform
	cfg.tokenSecret = secretString
	cfg.polkaKey = polkaKey

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
	serveMux.HandleFunc("PUT /api/users", cfg.updateUser)
	serveMux.HandleFunc("POST /api/login", cfg.loginUser)
	serveMux.HandleFunc("POST /api/refresh", cfg.refreshAccessToken)
	serveMux.HandleFunc("POST /api/revoke", cfg.revokeRefreshToken)

	// API Chirp Routes
	serveMux.HandleFunc("POST /api/chirps", cfg.createChirp)
	serveMux.HandleFunc("GET /api/chirps", cfg.getAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	serveMux.HandleFunc("POST /api/polka/webhooks", cfg.upgradeUserToChirpyRed)

	// Admin Routes
	serveMux.HandleFunc("GET /admin/metrics", cfg.returnMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.deleteAllUsers)

	err = server.ListenAndServe()

	if (err != nil) {
		log.Fatal(err)
	}

}


