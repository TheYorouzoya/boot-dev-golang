package main

import (
	"database/sql"
	"net/http"
	"time"
	"fmt"
	"github.com/TheYorouzoya/boot-dev-golang/Chirpy/internal/auth"
)


func (cfg *apiConfig) refreshAccessToken(writer http.ResponseWriter, request *http.Request) {
	headerToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Malformed authorization header", err)
		return
	}

	fetchedToken, err := cfg.dbQueries.GetRefreshToken(request.Context(), headerToken)
	if err != nil {
		if err == sql.ErrNoRows {
			responseError(writer, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		responseError(writer, http.StatusInternalServerError, "Error fetching token from DB", err)
		return
	}

	if fetchedToken.RevokedAt.Valid {
		err = fmt.Errorf("Given refresh token is revoked at %s", fetchedToken.RevokedAt.Time)
		responseError(writer, http.StatusUnauthorized, "Refresh token already revoked", err)
		return
	}

	if fetchedToken.ExpiresAt.Before(time.Now()) {
		responseError(writer, http.StatusUnauthorized, "Refresh token is expired", fmt.Errorf("expired refresh token"))
		return
	}


	accessToken, err := auth.MakeJWT(fetchedToken.UserID, cfg.tokenSecret, time.Duration(time.Hour * 24 * 60))
	if err != nil {
		responseError(writer, http.StatusInternalServerError, "Error generating new access token", err)
		return
	}

	responseJSON(writer, http.StatusOK, struct{Token string `json:"token"`}{Token: accessToken})
}


func (cfg *apiConfig) revokeRefreshToken(writer http.ResponseWriter, request *http.Request) {
	headerToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		responseError(writer, http.StatusUnauthorized, "Malformed authorization header", err)
		return
	}

	if _, err = cfg.dbQueries.RevokeRefreshToken(request.Context(), headerToken); err != nil {
		if err == sql.ErrNoRows {
			responseError(writer, http.StatusNotFound, "refresh token not found", err)
			return
		}
		responseError(writer, http.StatusInternalServerError, "error fetching refresh token", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
