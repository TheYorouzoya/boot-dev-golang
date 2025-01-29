package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/auth"
	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	Body 		string 		`json:"body"`
	UserID 		uuid.UUID 	`json:"user_id"`
}

func (cfg *apiConfig) createChirp(writer http.ResponseWriter, request *http.Request) {
	type chirpData struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Missing/Invalid authorization token in header", err)
		return
	}

	posterID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Invalid JWT token", err)
		return
	}

	decoder := json.NewDecoder(request.Body)
	requestData := chirpData{}
	if err := decoder.Decode(&requestData); err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON: %s", err), err)
		return
	}

	if len(requestData.Body) > 140 {
		responseError(writer, http.StatusInternalServerError, "Chirp is too long", nil)
		return
	}

	requestData.Body = cleanUpChirp(requestData.Body)
	// JWT determines the user posting the chirp
	requestData.UserID = posterID

	newChirp, err := cfg.dbQueries.CreateChirp(request.Context(), database.CreateChirpParams{
		Body: requestData.Body,
		UserID: requestData.UserID,})

	if err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %s", err), err)
		return
	}

	nChirp := Chirp(newChirp)
	responseJSON(writer, http.StatusCreated, nChirp)

}


func (cfg *apiConfig) getChirp(writer http.ResponseWriter, request *http.Request) {
	chirpID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Invalid UUID: %s", err), err)
	}

	chirpData, err := cfg.dbQueries.GetChirp(request.Context(), chirpID)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	chirp := Chirp(chirpData)
	responseJSON(writer, http.StatusOK, chirp)
}



func (cfg *apiConfig) getAllChirps(writer http.ResponseWriter, request *http.Request) {

	allChirps, err := cfg.dbQueries.AllChirps(request.Context())
	if err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error fetching chirps: %s", err), err)
	}

	chirps := make([]Chirp, len(allChirps))

	for i, chirp := range allChirps {
		chirps[i] = Chirp(chirp)
	}

	responseJSON(writer, http.StatusOK, chirps)
}


func cleanUpChirp(chirp string) string {
	words := strings.Split(chirp, " ")
	censorString := "****"
	badWords := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

	for index, word := range words {
		if badWords[strings.ToLower(word)] {
			// word is blacklisted
			words[index] = censorString
		}
	}

	return strings.Join(words, " ")
}
