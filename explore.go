package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	Id        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int `json:"chance"`
				ConditionValues []struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"condition_values"`
				MaxLevel int `json:"max_level"`
				Method   struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (cfg *apiConfig) commandExplore(args ...string) error {
	if len(args) == 0 {
		fmt.Println("Enter location area name to explore")
		return nil
	}

	locationAreaName := args[0]
	fmt.Printf("Exploring %s...\n", locationAreaName)
	fullURL, err := url.JoinPath(pokeApiBaseUrl, locationAreaEndpoint, locationAreaName)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	bytes, err := cfg.fetchWithCache(fullURL)
	if err != nil {
		return err
	}

	var locationArea LocationArea
	if err := json.Unmarshal(bytes, &locationArea); err != nil {
		return fmt.Errorf("error decoding location area: %w", err)
	}

	if len(locationArea.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, encounter := range locationArea.PokemonEncounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		}
	} else {
		fmt.Println("No Pokemon found in", locationAreaName)
	}

	return nil
}
