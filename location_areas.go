package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	locationAreaEndpoint = "/location-area"
	defaultLimit         = "20"
)

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

// commandMap prints the next set of location areas, using a cached URL or building a new request if needed.
func (cfg *apiConfig) commandMap(_ ...string) error {
	mapUrl := cfg.nextMapUrl
	if mapUrl == "" {
		parsedUrl, err := url.Parse(pokeApiBaseUrl + locationAreaEndpoint)
		if err != nil {
			return fmt.Errorf("error parsing URL: %v", err)
		}

		query := parsedUrl.Query()
		query.Set("offset", "0")
		query.Set("limit", defaultLimit) // if in future we want to make limit configurable, we can add a command argument for it and use that here instead of defaultLimit

		parsedUrl.RawQuery = query.Encode()
		mapUrl = parsedUrl.String()
	}

	normalizedUrl, err := normalizeUrl(mapUrl)
	if err != nil {
		return fmt.Errorf("error normalizing URL: %v", err)
	}

	cfg.fetchAndPrintLocationAreas(normalizedUrl)
	return nil
}

// commandMapBack prints the previous set of location areas, or notifies if already on the first page.
func (cfg *apiConfig) commandMapBack(_ ...string) error {
	if cfg.previousMapUrl == "" && cfg.nextMapUrl == "" {
		fmt.Println("Enter map to display the names of next 20 location areas of the Pokemon world")
		return nil
	} else if cfg.previousMapUrl == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	normalizedUrl, err := normalizeUrl(cfg.previousMapUrl)
	if err != nil {
		return fmt.Errorf("error normalizing URL: %v", err)
	}
	cfg.fetchAndPrintLocationAreas(normalizedUrl)
	return nil
}

// fetchAndPrintLocationAreas fetches location areas from the given URL if not cached, caches them, and prints their names.
// Given URL should be normalized to ensure consistent cache keys.
func (cfg *apiConfig) fetchAndPrintLocationAreas(normalizedMapUrl string) {
	bytes, isCached := cfg.mapCache.Get(normalizedMapUrl)
	if !isCached {
		res, err := http.Get(normalizedMapUrl)
		if err != nil {
			fmt.Printf("error fetching location areas: %v\n", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(res.Body)

		if res.StatusCode != http.StatusOK {
			fmt.Printf("unexpected status code: %d", res.StatusCode)
			return
		}

		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("error reading response body: %v", err)
			return
		}

		cfg.mapCache.Add(normalizedMapUrl, bytes)
	}

	var locationAreas LocationAreas
	if err := json.Unmarshal(bytes, &locationAreas); err != nil {
		fmt.Printf("error decoding location areas: %v", err)
		return
	}

	cfg.nextMapUrl = locationAreas.Next
	cfg.previousMapUrl = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}
}

// normalizeUrl parses a URL and re-encodes it with sorted query parameters
// This ensures consistent cache keys regardless of query parameter order
func normalizeUrl(urlString string) (string, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	parsedUrl.RawQuery = parsedUrl.Query().Encode() // Calling Query() and Encode() sorts the params lexicographically
	return parsedUrl.String(), nil
}
