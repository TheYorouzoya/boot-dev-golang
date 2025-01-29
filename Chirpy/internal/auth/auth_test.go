package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name 			string
		userID 			uuid.UUID
		secretString 	string
		expiresIn 		time.Duration
		wantErr			bool
	}{
		{
			name:			"valid token",
			userID:			uuid.New(),
			secretString:	"jwt-test-secret",
			expiresIn:		time.Hour,
			wantErr:		false,
		},
		{
			name:			"valid token short duration",
			userID:			uuid.New(),
			secretString:	"jwt-test-secret",
			expiresIn:		time.Second,
			wantErr:		false,
		},
		{
			name:			"empty secret",
			userID:			uuid.New(),
			secretString:	"",
			expiresIn:		time.Hour,
			wantErr:		true,
		},
		{
			name:         "zero duration",
			userID:       uuid.New(),
			secretString: "secret123",
			expiresIn:    0,
			wantErr:      true,
		},
		{
			name:         "negative duration",
			userID:       uuid.New(),
			secretString: "secret123",
			expiresIn:    -time.Hour,
			wantErr:      true,
		},
		{
			name:         "nil UUID",
			userID:       uuid.Nil,
			secretString: "secret123",
			expiresIn:    time.Hour,
			wantErr:      true,
		},

	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			token, err := MakeJWT(testCase.userID, testCase.secretString, testCase.expiresIn)
			if (err != nil) != testCase.wantErr {
				t.Errorf("MakeJWT() error = %v, wanted Err %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr {
				return
			}


			if len(token) <= 0 {
				t.Errorf("MakeJWT() returned an empty token")
				return
			}
			parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(tkn *jwt.Token) (interface{}, error) {
				return []byte(testCase.secretString), nil
			})
			if err != nil {
				t.Errorf("Could not parse the token returned from MakeJWT(): %v", err)
				return
			}

			claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
			if !ok {
				t.Errorf("Could not get claims for the returned token")
				return
			}

			if claims.Issuer != "chirpy" {
				t.Errorf("Expected issuer to be 'chirpy', got %v", claims.Issuer)
			}
			if claims.Subject != testCase.userID.String() {
				t.Errorf("Expected subject to be %v, got %v", testCase.userID, claims.Subject)
			}
			expectedExp := time.Now().Add(testCase.expiresIn)
			epsilon := time.Second

			if claims.ExpiresAt == nil {
				t.Error("ExpiresAt claim is nil")
			} else {
				diff := claims.ExpiresAt.Time.Sub(expectedExp)
				if diff < -epsilon || diff > epsilon {
					t.Errorf("ExpiresAt above tolerance, expected: %v, got: %v", expectedExp, claims.ExpiresAt.Time)
				}
			}

		})
	}
}


func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name 			string
		setupToken 		func() string
		secretString 	string
		wantUserID 		uuid.UUID
		wantErr 		bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				userID := uuid.New()
				token, _ := MakeJWT(userID, "jwt-test-secret", time.Hour)
				return token
			},
			secretString: "jwt-test-secret",
			wantErr: false,
		},
		{
            name: "expired token",
            setupToken: func() string {
                userID := uuid.New()
                token, _ := MakeJWT(userID, "jwt-test-secret", -time.Hour)
                return token
            },
            secretString: "jwt-test-secret",
            wantErr:     true,
        },
        {
            name: "wrong secret",
            setupToken: func() string {
                userID := uuid.New()
                token, _ := MakeJWT(userID, "jwt-wrong-secret", time.Hour)
                return token
            },
            secretString: "jwt-test-secret",
            wantErr:     true,
        },
        {
            name:         "invalid token format",
            setupToken:   func() string { return "not.a.token" },
            secretString: "jwt-test-secret",
            wantErr:     true,
        },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token := test.setupToken()
			userID, err := ValidateJWT(token, test.secretString)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateJWT() returned error: %v, expected error: %v", err, test.wantErr)
				return
			}
			if !test.wantErr && userID == uuid.Nil {
				t.Error("ValidateJWT() returned nil UUID for a valid token")
				return
			}
			if test.wantErr && userID != uuid.Nil {
				t.Error("ValidateJWT() returned a non-nil UUID for a valid token")
			}
		})
	}

}


func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name 			string
		testHeader 		http.Header
		expectedToken 	string
		wantErr 		bool
	}{
		{
			name:			"valid header",
			testHeader:		http.Header{"Authorization": []string{"Bearer testTokenString"}},
			expectedToken:	"testTokenString",
			wantErr:		false,
		},
		{
			name: 			"Missing token header",
			testHeader: 	http.Header{"Authorization": []string{},},
			expectedToken: 	"",
			wantErr:       	true,
		},
		{
			name:			"Header contains more than one bearer token string",
			testHeader:		http.Header{"Authorization": []string{"Bearer testTokenString another"}},
			expectedToken:	"",
			wantErr:		true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := GetBearerToken(test.testHeader)
			if (err != nil) != test.wantErr {
				t.Errorf("GetBearerToken() returned error: %v, expected error: %v", err, test.wantErr)
				return
			}
			if test.expectedToken != token {
				t.Errorf("token mismatch, wanted: %s, got: %s", test.expectedToken, token)
			}
		})
	}

}
