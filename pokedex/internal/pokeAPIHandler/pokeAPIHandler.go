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
		Name *string `json:"name"`
		URL  *string `json:"url"`
	} `json:"results"`
}


type PokeAPIExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon		struct {
			Name	*string		`json:"name"`
			URL		*string		`json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}


type PokeAPIPokemonResponse struct {
	Name				*string		`json:"name"`
	BaseExperience		int			`json:"base_experience"`
	Height				int			`json:"height"`
	Weight				int			`json:"weight"`
	Stats				[]struct{
		BaseValue		int			`json:"base_stat"`
		Stat			struct {
			Name		*string			`json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types				[]struct{
		Type 			struct{
			Name		*string			`json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// Query the PokeAPI for a given query and return the raw byte response
// will return cached response if query is in cache
func QueryAPI(query *string, cache *pokecache.Cache, isCached bool) ([]byte, error) {
	if query == nil {
		return nil, fmt.Errorf("Query is empty")
	}

	cacheResponse, ok := cache.Get(query)

	if (ok) {
		return cacheResponse, nil
	}

	res, err := http.Get(*query)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("API response failed with status code: %d\n", res.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	if isCached {
		cache.Add(query, body)
	}

	return body, nil
}

func QueryMap(query *string, cache *pokecache.Cache) (PokeAPIMapResponse, error) {
	var apiResponse PokeAPIMapResponse

	response, err := QueryAPI(query, cache, true)

	if err != nil {
		return PokeAPIMapResponse{}, err
	}

	err = json.Unmarshal(response, &apiResponse)
	if err != nil {
		return PokeAPIMapResponse{}, err
	}

	return apiResponse, nil
}


func QueryExplore(query *string, cache *pokecache.Cache) (PokeAPIExploreResponse, error) {
	var apiResponse PokeAPIExploreResponse

	response, err := QueryAPI(query, cache, true)

	if err != nil {
		return PokeAPIExploreResponse{}, err
	}

	err = json.Unmarshal(response, &apiResponse)

	if err != nil {
		return PokeAPIExploreResponse{}, err
	}

	return apiResponse, nil
}


func QueryPokemon(query *string, cache *pokecache.Cache) (PokeAPIPokemonResponse, error) {
	var apiResponse PokeAPIPokemonResponse

	response, err := QueryAPI(query, cache, false)

	if err != nil {
		return PokeAPIPokemonResponse{}, err
	}

	err = json.Unmarshal(response, &apiResponse)

	if err != nil {
		return PokeAPIPokemonResponse{}, err
	}

	return apiResponse, nil
}
