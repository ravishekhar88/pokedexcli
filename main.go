package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ravishekhar88/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

var (
	errExitRequested = errors.New("exit requested")
	cmds             map[string]cliCommand
)

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
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    cfg.commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect an already caught Pokemon",
			callback:    cfg.commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the names of all the caught Pokemon",
			callback:    cfg.commandPokedex,
		},
	}
}

func main() {
	logFile, err := initLog()
	if err != nil {
		fmt.Printf("Could not initialize logging: %v\n", err)
		return
	}
	defer func() {
		if closeErr := logFile.Close(); closeErr != nil {
			fmt.Printf("Could not close log file: %v\n", closeErr)
		}
	}()

	dbPath, err := resolveDBPath()
	if err != nil {
		log.Printf("Could not resolve database path: %v", err)
		fmt.Println("Could not resolve database path, check cli.log for details")
		return
	}

	db, err := initDB(dbPath)
	if err != nil {
		log.Printf("Could not initialize database: %v", err)
		fmt.Println("Could not initialize database, check cli.log for details")
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("Could not close database: %v", err)
		}
	}()

	pokemons, err := loadPokemons(db)
	if err != nil {
		log.Printf("Could not load pokemons from database: %v", err)
		fmt.Println("Could not load pokemons from database, check cli.log for details")
		return
	}

	cfg := apiConfig{
		db:           db,
		pokeApiCache: pokecache.NewCache(10 * time.Second),
		pokemons:     pokemons,
	}

	initializeCommands(cfg)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		input, ok := scanInput(scanner)
		if !ok {
			break
		}
		if shouldExit := handleInput(input); shouldExit {
			return
		}
	}
}

func scanInput(scanner *bufio.Scanner) (string, bool) {
	fmt.Print("Pokedex > ")
	if !scanner.Scan() {
		if scanErr := scanner.Err(); scanErr != nil {
			log.Printf("input scanner error: %v", scanErr)
		}
		return "", false
	}

	return strings.TrimSpace(scanner.Text()), true
}

func handleInput(input string) bool {
	cleanedInput := cleanInput(input)
	if len(cleanedInput) == 0 {
		return false
	}

	userCmd := cleanedInput[0]
	args := cleanedInput[1:]

	cmd, found := cmds[userCmd]
	if !found {
		fmt.Println("Unknown command")
		return false
	}

	err := cmd.callback(args...)
	if err == nil {
		return false
	}
	if errors.Is(err, errExitRequested) {
		return true
	}

	log.Printf("Error executing command: %s: %v\n", input, err)
	fmt.Printf("Error executing command: %s\n", input)
	return false
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
	return errExitRequested
}

func initLog() (*os.File, error) {
	f, err := os.OpenFile("cli.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return f, nil
}

func resolveDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not find user config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "pokedexcli")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("could not create app config directory: %w", err)
	}

	return filepath.Join(appDir, "pokedex.db"), nil
}

func cleanInput(text string) []string {
	fields := strings.Fields(text)
	for i, field := range fields {
		fields[i] = strings.ToLower(field)
	}
	return fields
}
