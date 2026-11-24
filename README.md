# Go REPL Pokedex

An interactive command-line Pokedex written in Go. It uses the [PokeAPI](https://pokeapi.co/) to browse location areas, explore Pokémon encounters, and attempt to catch Pokémon.

## Requirements

- Go 1.27 or newer
- Internet access to reach PokeAPI

## Run

From the project directory:

```sh
go run .
```

Or build and run a binary:

```sh
go build -o pokedex .
./pokedex
```

## Commands

| Command | Description |
| --- | --- |
| `help` | Show available commands |
| `map` | List 20 location areas and advance to the next page |
| `mapb` | Return to the previous page of location areas |
| `explore <location-area>` | List Pokémon encountered in a location area |
| `catch <pokemon>` | Attempt to catch a Pokémon |
| `inspect <pokemon>` | Show details for a Pokémon caught during this session |
| `save [filename]` | Save caught Pokémon to a file in `saves/`; defaults to `save_1.gob` |
| `exit` | Close the Pokedex |

For example:

```text
Pokedex > map
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > exit
```

## Saving and loading

At startup, the Pokedex attempts to load `saves/save_1.gob`. Use `save` to write the current Pokédex to that file, or `save <filename>` to save it under a different filename in `saves/`. Save files use Go's gob encoding.

Loading is automatic and currently only uses the default `save_1.gob` file; there is no `load` command. A custom-named save file is not loaded automatically. Pokémon names are case-insensitive when catching.

## Future improvements

- Add a `load [filename]` command to the REPL command configuration, allowing the user to load a chosen save and displaying the name of the save file that was loaded.

## Network and caching

API requests use a 10-second client timeout. Successfully decoded responses are cached in memory. Reads refresh each entry's age, and a background reaper scans every 10 minutes to remove entries older than 10 minutes.

## Tests

Run the test suite with:

```sh
go test ./...
```
