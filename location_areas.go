package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const defaultLimit = "20"

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

	err = cfg.fetchAndPrintLocationAreas(normalizedUrl)
	if err != nil {
		return err
	}

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

	err = cfg.fetchAndPrintLocationAreas(normalizedUrl)
	if err != nil {
		return err
	}

	return nil
}

// fetchAndPrintLocationAreas fetches location areas from the given URL if not cached, caches them, and prints their names.
// Given URL should be normalized to ensure consistent cache keys.
func (cfg *apiConfig) fetchAndPrintLocationAreas(normalizedMapUrl string) error {
	bytes, err := cfg.fetchWithCache(normalizedMapUrl)
	if err != nil {
		return err
	}

	var locationAreas LocationAreas
	if err := json.Unmarshal(bytes, &locationAreas); err != nil {
		return fmt.Errorf("error decoding location areas: %w", err)
	}

	cfg.nextMapUrl = locationAreas.Next
	cfg.previousMapUrl = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
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
