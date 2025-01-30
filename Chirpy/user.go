package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/auth"
	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID 				uuid.UUID 	`json:"id"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	Email 			string 		`json:"email"`
	HashedPassword 	string 		`json:"-"`
	IsChirpyRed		bool 		`json:"is_chirpy_red"`
}

type userData struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		User
	}



func (cfg *apiConfig) createUser(writer http.ResponseWriter, request *http.Request) {

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

	decoder := json.NewDecoder(request.Body)
	uData := userData{}

	if err := decoder.Decode(&uData); err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON: %s", err), err)
		return
	}

	defaultExpirationTime := time.Hour

	usrData, err := cfg.dbQueries.GetUserWithEmail(request.Context(), uData.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			responseError(writer, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}
		responseError(writer, http.StatusInternalServerError, "Error fetching user from DB", err)
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


func (cfg *apiConfig) updateUser(writer http.ResponseWriter, request *http.Request) {

	accessToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Missing/Malformed auth token in header", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Invalid auth token", err)
		return
	}

	decoder := json.NewDecoder(request.Body)
	uData := userData{}

	if err := decoder.Decode(&uData); err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON: %s", err), err)
		return
	}

	hashedPassword, err := auth.HashPassword(uData.Password)
	if err != nil {
		responseError(writer, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUser(request.Context(), database.UpdateUserParams{
		Email: uData.Email,
		HashedPassword: hashedPassword,
		ID: userID,
	})
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			responseError(writer, http.StatusBadRequest, "Email already taken", err)
			return
		}
		responseError(writer, http.StatusInternalServerError, "Error updating user in DB", err)
		return
	}

	finalResponse := struct{
		Email string `json:"email"`
		ID uuid.UUID `json:"id"`
	}{
		Email: updatedUser.Email,
		ID: updatedUser.ID,
	}

	responseJSON(writer, http.StatusOK, finalResponse)
}


func (cfg *apiConfig) upgradeUserToChirpyRed(writer http.ResponseWriter, request *http.Request) {
	type requestData struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Malformed/Missing API key", err)
		return
	}

	if cfg.polkaKey != apiKey {
		responseError(writer, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	var data requestData

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&data); err != nil {
		responseError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	if data.Event != "user.upgraded" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(data.Data.UserID)
	if err != nil {
		responseError(writer, http.StatusBadRequest, "Malformed UUID", err)
		return
	}

	_, err = cfg.dbQueries.UpgradeUserToChirpyRed(request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			responseError(writer, http.StatusNotFound, "User does not exist", err)
			return
		}
		responseError(writer, http.StatusInternalServerError, "Error upgrading user", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
