package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	pokego "github.com/vinagoros/go-repl/internal"
)

// passar endpoint para módulo base de get
var baseEndpoint string = "https://pokeapi.co/api/v2/"

type Pokedex struct {
	Name string
}

type location struct {
	Name string `json:"name"`
	url  string `json:"-"` //unexported
}

type locationAreas struct {
	Count    int        `json:"-"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

type locationArea struct {
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type pokemon struct {
	Name string `json:"name"`
}

type pokemonDetails struct {
	Name           string `json:"name"`
	BaseExperience int32  `json:"base_experience"`
	Height         int32  `json:"height"`
	Nr             int32  `json:"id"`
	Weight         int32  `json:"weight"`
	Stats          []struct {
		BaseStatValue int32 `json:"base_stat"`
		Stat          struct {
			StatName string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type config struct {
	availableCommands map[string]cliCommand
	next              string
	previous          string
}

type cliCommand struct {
	name         string
	description  string
	callbackFunc func(*config, string) error
}

func exit(c *config, arg string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func pokeMap(c *config, arg string) error {
	endpoint := baseEndpoint + "location-area/" //passar endpoint para modulo base
	if c.next == "" && c.previous != "" {
		return fmt.Errorf("No further locations available")
	}
	if c.next != "" {
		endpoint = c.next
	}
	resp, err := pokego.PokeGet[locationAreas](fmt.Sprintf("%s/?limit=%d", endpoint, 20))
	if err != nil {
		return err
	}
	c.next = resp.Next
	c.previous = resp.Previous
	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func pokeMapB(c *config, arg string) error {
	var endpoint string
	if c.previous != "" {
		endpoint = c.previous
	} else {
		fmt.Println("you're on the first page")
		return errors.New("Cannot go back further")
	}
	resp, err := pokego.PokeGet[locationAreas](endpoint)
	c.next = resp.Next
	c.previous = resp.Previous
	for _, result := range resp.Results {
		fmt.Println(result.Name)
	}

	return err
}

func explore(c *config, locationAreaToQuery string) error {
	endpoint := baseEndpoint + "location-area/"
	if len(locationAreaToQuery) < 1 {
		return fmt.Errorf("Error exploring: No area provided")
	}
	resp, err := pokego.PokeGet[locationArea](fmt.Sprintf("%s/%s", endpoint, locationAreaToQuery))
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve Pokemon data for location %s: %w", locationAreaToQuery, err)
	}
	for _, encounter := range resp.PokemonEncounters {
		fmt.Println("- ", encounter.Pokemon.Name)
	}
	return nil
}

func catch(c *config, pokemonName string) error {
	var maxBaseExperience int = 324
	if len(pokemonName) <= 1 {
		return fmt.Errorf("Error trying to catch: Invalid pokemon name")
	}
	fmt.Println("Throwing a Pokeball at ", pokemonName)
	endpoint := baseEndpoint + fmt.Sprintf("pokemon/%s", strings.ToLower(pokemonName))
	fmt.Println("url:", endpoint)
	resp, err := pokego.PokeGet[pokemonDetails](endpoint)
	if err != nil {
		return fmt.Errorf("Error while trying to retrieve Pokemon %s: %w", pokemonName, err)
	}
	catchProb := rand.Int31n(int32(maxBaseExperience))
	if catchProb >= resp.BaseExperience {
		fmt.Println("You've catched", resp.Name, "!")
	}

	return nil
}

func help(c *config, arg string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for key, supportedCommand := range supportedCommands {
		fmt.Printf("%s: %s\n", key, supportedCommand.description)
	}
	return nil
}

var supportedCommands = map[string]cliCommand{
	"exit": {
		name:         "exit",
		description:  "Exit the Pokedex",
		callbackFunc: exit,
	},
	"map": {
		name:         "map",
		description:  "Retrieves 20 location areas of the Pokemon World",
		callbackFunc: pokeMap,
	},
	"mapb": {
		name:         "mapb",
		description:  "Retrieves previous 20 location ares of the Pokemon World",
		callbackFunc: pokeMapB,
	},
	"explore": {
		name:         "explore",
		description:  "Retrieves all Pokemon encounters in the area",
		callbackFunc: explore,
	},
	"catch": {
		name:         "catch",
		description:  "Attempts to catch a specified Pokemon",
		callbackFunc: catch,
	},
}

func main() {
	supportedCommands["help"] = cliCommand{
		name:         "help",
		description:  "Displays a help message",
		callbackFunc: help,
	}
	conf := config{
		availableCommands: supportedCommands,
	}
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		result := reader.Text()
		firstWord, arg, _ := strings.Cut(result, " ")
		command := strings.Trim(strings.ToLower(firstWord), " ")
		for key, supportCommand := range supportedCommands {
			if key == command {
				err := supportCommand.callbackFunc(&conf, arg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
