package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}


func main() {
	serveMux := http.NewServeMux()
	var server http.Server
	var apiCfg apiConfig

	server.Addr = ":8080"
	server.Handler = serveMux

	serveMux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix(
				"/app",http.FileServer(http.Dir(".")))))

	serveMux.HandleFunc("GET /api/healthz", readinessCheck)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.returnMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirp)
	err := server.ListenAndServe()

	if (err != nil) {
		log.Fatal(err)
	}

}


