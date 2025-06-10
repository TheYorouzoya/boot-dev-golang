package pokeAPIHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/TheYorouzoya/boot-dev-golang/pokedex/internal/pokecache"
)

type PokeAPIMapResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func QueryMap(query *string, cache *pokecache.Cache) (PokeAPIMapResponse, error) {
	if query == nil {
		return PokeAPIMapResponse{}, fmt.Errorf("Query is empty")
	}

	var apiResponse PokeAPIMapResponse

	cacheResponse, ok := cache.Get(query)

	if (ok) {
		err := json.Unmarshal(cacheResponse, &apiResponse)
		if err != nil {
			return PokeAPIMapResponse{}, err
		}
		return apiResponse, nil
	}

	// make the api request to get map data
	res, err := http.Get(*query)
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

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return PokeAPIMapResponse{}, err
	}

	cache.Add(query, body)

	return apiResponse, nil
}
