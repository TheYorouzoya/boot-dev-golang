package main

import (
	"fmt"
	"net/http"
)


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, request)
	})
}


func (cfg *apiConfig) returnMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
		cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetMetrics(writer http.ResponseWriter, request *http.Request) {
	// reset fileserver hits to 0
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(http.StatusOK)
}
