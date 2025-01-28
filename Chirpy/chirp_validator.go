package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)


func validateChirp(writer http.ResponseWriter, request *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type successResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(request.Body)
	chirpData := chirp{}
	if err := decoder.Decode(&chirpData); err != nil {
		errResponse := errorResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)}
		responseJSON(writer, http.StatusInternalServerError, errResponse)
		return
	}

	if len(chirpData.Body) > 140 {
		errResponse := errorResponse{Error: "Chirp is too long"}
		responseJSON(writer, http.StatusBadRequest, errResponse)
		return
	}

	succResponse := successResponse{Valid: true}
	responseJSON(writer, http.StatusOK, succResponse)
	return
}
