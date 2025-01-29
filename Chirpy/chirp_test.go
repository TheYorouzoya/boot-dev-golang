package main

import (
    "testing"

)


// func TestValidateChirpHandler(t *testing.T) {
//     tests := []struct {
//         name           string
//         inputJSON      string
//         expectedStatus int
//         expectedBody   string
//     }{
//         {
//             name:           "valid chirp no profanity",
//             inputJSON:      `{"body": "Hello, World!"}`,
//             expectedStatus: http.StatusOK,
//             expectedBody:   `{"cleaned_body":"Hello, World!"}`,
//         },
//         {
//             name:           "valid chirp with profanity",
//             inputJSON:      `{"body": "Hello kerfuffle!"}`,
//             expectedStatus: http.StatusOK,
//             expectedBody:   `{"cleaned_body":"Hello kerfuffle!"}`,
//         },
//         {
//             name:           "chirp too long",
//             inputJSON:      `{"body": "` + strings.Repeat("x", 141) + `"}`,
//             expectedStatus: http.StatusBadRequest,
//             expectedBody:   `{"error":"Chirp is too long"}`,
//         },
//         {
// 			name:           "invalid JSON",
// 			inputJSON:      `{"body": not valid json}`,
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   `{"error":"Error decoding JSON: invalid character 'o' in literal null (expecting 'u')"}`,
// 		},
//     }

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest(
// 				http.MethodPost,
// 				"/api/validate_chirp",
// 				strings.NewReader(tt.inputJSON),
// 			)
// 			req.Header.Set("Content-Type", "application/json")
// 			w := httptest.NewRecorder()

// 			validateChirp(w, req)

// 			if w.Code != tt.expectedStatus {
// 				t.Errorf("wanted status %v, got %v", tt.expectedStatus, w.Code)
// 			}

// 			gotBody := strings.TrimSpace(w.Body.String())
// 			if gotBody != tt.expectedBody {
// 				t.Errorf("wanted response body %v, got %v", tt.expectedBody, gotBody)
// 			}
// 		})
// 	}
// }


func TestCleanUpChirp(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "no profanity",
            input:    "hello world",
            expected: "hello world",
        },
        {
            name:     "single profane word",
            input:    "hello kerfuffle world",
            expected: "hello **** world",
        },
        {
            name:     "multiple profane words",
            input:    "kerfuffle hello sharbert",
            expected: "**** hello ****",
        },
        {
			name:     "profane word with punctuation",
			input:	  "hello world kerfuffle!",
			expected: "hello world kerfuffle!",
		},
		{
			name:     "profane word with capitalization",
			input:	  "hello world KerFuFfle ",
			expected: "hello world **** ",
		},
		{
			name:     "empty string",
			input:	  "",
			expected: "",
		},
		{
			name:     "single word",
			input:	  "kerfuffle",
			expected: "****",
		},
		{
			name:     "multiple consecutive spaces",
			input:    "hello    kerfuffle     world",
			expected: "hello    ****     world",
		},
		{
			name:     "word containing profane word as substring",
			input:    "kerfuffles sharbertify fornaxious",
			expected: "kerfuffles sharbertify fornaxious",
		},
		{
			name:     "mixed spaces and punctuation",
			input:    "sharbert!   kerfuffle?    fornax...",
			expected: "sharbert!   kerfuffle?    fornax...",
		},
		{
			name:     "multiple punctuation marks",
			input:    "sharbert!!! kerfuffle??? fornax...",
			expected: "sharbert!!! kerfuffle??? fornax...",
		},
	}

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := cleanUpChirp(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
