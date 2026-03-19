package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ravishekhar88/pokedexcli/internal/pokecache"
)

const (
	pokeApiBaseUrl       = "https://pokeapi.co/api/v2"
	locationAreaEndpoint = "/location-area"
	pokemonEndpoint      = "/pokemon"
)

type apiConfig struct {
	pokeApiCache   *pokecache.Cache
	nextMapUrl     string
	previousMapUrl string
	pokemons       map[string]Pokemon
}

// fetchWithCache retrieves data from the cache if available, otherwise fetches it via HTTP, and adds the response to the cache.
func (cfg *apiConfig) fetchWithCache(fullURL string) ([]byte, error) {
	bytes, isCached := cfg.pokeApiCache.Get(fullURL)
	if !isCached {
		res, err := http.Get(fullURL)
		if err != nil {
			return nil, fmt.Errorf("error fetching location area details: %w", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(res.Body)

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		cfg.pokeApiCache.Add(fullURL, bytes)
	}

	return bytes, nil
}
