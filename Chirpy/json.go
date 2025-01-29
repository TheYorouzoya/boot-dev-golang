package main

import (
	"log"
	"net/http"
	"encoding/json"
)

type errorResponse struct {
		Error string `json:"error"`
	}

func responseError(writer http.ResponseWriter, status int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	responseJSON(writer, status, errorResponse{
		Error: msg,
	})
}

func responseJSON(writer http.ResponseWriter, status int, rawData interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(rawData)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(status)
	writer.Write(dat)
}
