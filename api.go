package main

import (
	"database/sql"
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
	db             *sql.DB
	pokeApiCache   *pokecache.Cache
	nextMapUrl     string
	previousMapUrl string
	pokemons       map[string]Pokemon
}

// fetchWithCache retrieves data from the cache if available, otherwise fetches it via HTTP, and adds the response to the cache.
func (cfg *apiConfig) fetchWithCache(fullURL string, resourceType string) ([]byte, error) {
	bytes, isCached := cfg.pokeApiCache.Get(fullURL)
	if !isCached {
		res, err := http.Get(fullURL)
		if err != nil {
			return nil, fmt.Errorf("error fetching %s details: %w", resourceType, err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing %s response body: %v", resourceType, err)
			}
		}(res.Body)

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code for %s: %d", resourceType, res.StatusCode)
		}

		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading %s response body: %w", resourceType, err)
		}

		cfg.pokeApiCache.Add(fullURL, bytes)
	}

	return bytes, nil
}
