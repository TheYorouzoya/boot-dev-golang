package pokeAPIHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeAPIMapResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func QueryMap(query string) (PokeAPIMapResponse, error) {
	// make the api request to get map data
	res, err := http.Get(query)
	if err != nil {
		return PokeAPIMapResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return PokeAPIMapResponse{}, fmt.Errorf("API response failed with status code: %d\n", res.StatusCode)
	}

	if err != nil {
		return PokeAPIMapResponse{}, err
	}

	var apiResponse PokeAPIMapResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return PokeAPIMapResponse{}, err
	}

	return apiResponse, nil
}
