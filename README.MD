# Pokedex CLI

A command-line interface tool for exploring the Pokemon world using the PokeAPI. Navigate through location areas and discover Pokemon in different regions.

## Features

- Browse through Pokemon world location areas
- Navigate forward and backward through paginated results
- Explore specific locations to find Pokemon
- Built-in caching system for improved performance
- Simple, interactive REPL interface

## Installation

```bash
go get github.com/ravishekhar88/pokedexcli
```

Or clone and build locally:

```bash
git clone https://github.com/ravishekhar88/pokedexcli.git
cd pokedexcli
go build
```

## Usage

Run the application:

```bash
go run main.go
```

Or if you've built the binary:

```bash
./pokedexcli
```

### Available Commands

- `help` - Displays the help message with all available commands
- `exit` - Exit the Pokedex application
- `map` - Displays the names of next 20 location areas in the Pokemon world
- `mapb` - Displays the names of previous 20 location areas in the Pokemon world
- `explore <location-area>` - Displays the names of all Pokemon in a specific location area

### Example Session

```
Pokedex > help
Welcome to the Pokedex!
Usage:

help: Displays the help message
exit: Exit the Pokedex
map: Displays the names of next 20 location areas of the Pokemon world
mapb: Displays the names of previous 20 location areas of the Pokemon world
explore: Displays the names of all the Pokemon in a location area

Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > mapb
you're on the first page

Pokedex > explore canalave-city-area
tentacool
tentacruel
...

Pokedex > exit
Closing the Pokedex... Goodbye!
```

## Project Structure

```
pokedexcli/
├── main.go                          # Main application entry point and REPL
├── location_areas.go                # Location area fetching and navigation
├── explore.go                       # Pokemon exploration in location areas
├── internal/
│   └── pokecache/
│       ├── pokecache.go            # Caching implementation
│       └── pokecache_test.go       # Cache tests
├── go.mod                          # Go module definition
└── README.MD                       # This file
```

## Technical Details

### Caching

The application implements an in-memory cache with automatic cleanup:
- Cache entries expire after 10 seconds
- Thread-safe implementation using sync.RWMutex
- Automatic background reaping of expired entries

### API

This project uses the [PokeAPI](https://pokeapi.co/) to fetch Pokemon data.

## Requirements

- Go 1.25.1 or higher

## License

MIT License

Copyright (c) 2026

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
