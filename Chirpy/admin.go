package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) deleteAllUsers(writer http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	if err := cfg.dbQueries.DeleteAllUsers(request.Context()); err != nil {
		responseError(writer, http.StatusInternalServerError, fmt.Sprintf("Error deleting data: %s", err), err)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
