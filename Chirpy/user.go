package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/auth"
	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID 				uuid.UUID 	`json:"id"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	Email 			string 		`json:"email"`
	HashedPassword 	string 		`json:"-"`
}



func (cfg *apiConfig) createUser(writer http.ResponseWriter, request *http.Request) {
	type userData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(request.Body)
	uData := userData{}

	if err := decoder.Decode(&uData); err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON: %s", err), err)
		return
	}

	passHash, err := auth.HashPassword(uData.Password)
	if err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %s", err), err)
	}

	newUser, err := cfg.dbQueries.CreateUser(request.Context(), database.CreateUserParams{
		Email: uData.Email,
		HashedPassword: passHash,
	})

	if err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %s", err), err)
		return
	}

	nUser := User(newUser)
	responseJSON(writer, http.StatusCreated, nUser)
}


func (cfg *apiConfig) loginUser(writer http.ResponseWriter, request *http.Request) {
	type userData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type tokenResponse struct {
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		User
	}

	decoder := json.NewDecoder(request.Body)
	uData := userData{}

	if err := decoder.Decode(&uData); err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON: %s", err), err)
		return
	}

	defaultExpirationTime := time.Hour

	usrData, err := cfg.dbQueries.GetUser(request.Context(), uData.Email)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	user := User(usrData)

	if err = auth.CheckPasswordHash(uData.Password, user.HashedPassword); err != nil {
		responseError(writer, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, defaultExpirationTime)
	if err != nil {
		responseError(writer, http.StatusInternalServerError, "Error creating JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		responseError(writer, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	_, err = cfg.dbQueries.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		responseError(writer, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}


	finalResponse := tokenResponse{
		Token: accessToken,
		RefreshToken: refreshToken,
		User: user,
	}

	responseJSON(writer, http.StatusOK, finalResponse)
}
