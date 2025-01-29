package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}


func CheckPasswordHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	if len(tokenSecret) <= 0 {
		return "", fmt.Errorf("secret string cannot be empty")
	}

	if expiresIn <= 0 {
		return "", fmt.Errorf("expiresIn duration must be positive")
	}

	if userID == uuid.Nil {
		return "", fmt.Errorf("given user UUID cannot be nil")
	}

	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedStringToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedStringToken, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	if len(tokenSecret) <= 0 {
		return uuid.Nil, fmt.Errorf("secret string cannot be empty")
	}
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	} else if parsedClaims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if parsedClaims.ExpiresAt.Time.Before(time.Now()) {
			return uuid.Nil, fmt.Errorf("given JWT token has expired")
		}
		uid, parseErr := uuid.Parse(parsedClaims.Subject)
		if parseErr != nil {
			return uuid.Nil, fmt.Errorf("invalid uuid, cannot proceed")
		}
		return uid, nil
	} else {
		return uuid.Nil, fmt.Errorf("unknown claims type, cannot proceed")
	}
}


func GetBearerToken(headers http.Header) (string, error) {
	tokenHeader := headers.Get("Authorization")

	if tokenHeader == "" {
		return tokenHeader, fmt.Errorf("Authorization header doesn't exist")
	}

	tokenFields := strings.Fields(tokenHeader)

	if len(tokenFields) != 2 || tokenFields[0] != "Bearer" {
		return "", fmt.Errorf("Malformed authorization header")
	}

	return tokenFields[1], nil;
}


func MakeRefreshToken() (string, error) {
	tokenLength := 32
	tokenBytes := make([]byte, tokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes)
	return token, nil
}
