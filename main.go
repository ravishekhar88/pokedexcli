package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ravishekhar88/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

var cmds map[string]cliCommand

func initializeCommands(cfg apiConfig) {
	cmds = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays the help message",
			callback:    cfg.commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    cfg.commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of next 20 location areas of the Pokemon world",
			callback:    cfg.commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas of the Pokemon world",
			callback:    cfg.commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of all the Pokemon in a location area",
			callback:    cfg.commandExplore,
		},
	}
}

func main() {
	cfg := apiConfig{
		pokeApiCache: pokecache.NewCache(10 * time.Second),
	}

	initLog()
	initializeCommands(cfg)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		cleanedInput := cleanInput(input)
		userCmd := cleanedInput[0]
		args := cleanedInput[1:]

		cmd, found := cmds[userCmd]
		if !found {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(args...)
		if err != nil {
			log.Printf("Error executing command %s: %v\n", userCmd, err)
			fmt.Printf("Error executing command %s: \n", userCmd)
		}
	}
}

func (cfg *apiConfig) commandHelp(...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for key, cmd := range cmds {
		fmt.Printf("%s: %s\n", key, cmd.description)
	}

	return nil
}

func (cfg *apiConfig) commandExit(...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func initLog() {
	f, err := os.OpenFile("cli.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not open log file: %v\n", err)
		os.Exit(1)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Could not close log file: %v\n", err)
		}
	}(f)

	log.SetOutput(f)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func cleanInput(text string) []string {
	fields := strings.Fields(text)
	for i, field := range fields {
		fields[i] = strings.ToLower(field)
	}
	return fields
}
