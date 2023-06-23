package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Nikhil-Kumar21/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient       pokeapi.Client
	nextLocationAreaURL *string
	prevLocationAreaURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
}

type clicommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Printf("Welcome to the Pokedex!\n\n")
	fmt.Printf("Usage:\n\n")

	commands := getCommands()

	for _, cmd := range commands {
		fmt.Printf("%s - %s\n", cmd.name, cmd.description)
	}

	fmt.Println("")
	return nil
}

func commandExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args ...string) error {

	resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationAreaURL)

	if err != nil {
		return err
	}

	fmt.Println("Location areas:")

	for _, ar := range resp.Results {
		fmt.Printf(" - %s\n", ar.Name)
	}
	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous

	return nil
}
func commandMapb(cfg *config, args ...string) error {

	if cfg.prevLocationAreaURL == nil {
		return errors.New("you are on the first page")
	}
	resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationAreaURL)

	if err != nil {
		return err
	}

	fmt.Println("Location areas:")

	for _, ar := range resp.Results {
		fmt.Printf(" - %s\n", ar.Name)
	}
	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no location area provided")
	}

	locationAreaName := args[0]
	locationArea, err := cfg.pokeapiClient.GetLocationArea(locationAreaName)

	if err != nil {
		return err
	}

	fmt.Printf("Pokemon in %s:\n", locationArea.Name)

	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no pokemon name provided")
	}

	pokemonName := args[0]
	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)

	if err != nil {
		return err
	}

	randNum := rand.Intn(pokemon.BaseExperience)
	threshold := 50

	fmt.Println(pokemon.BaseExperience, randNum, threshold)

	if randNum > threshold {
		return fmt.Errorf("failed to catch %s", pokemonName)
	}

	fmt.Printf("%s was caught\n", pokemonName)
	cfg.caughtPokemon[pokemonName] = pokemon
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no pokemon name provided")
	}

	pokemonName := args[0]
	pokemon, ok := cfg.caughtPokemon[pokemonName]

	if !ok {
		return errors.New("pokemon not caught till now")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf(" - %s\n", typ.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Pokemon in Pokedex:")
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func getCommands() map[string]clicommand {
	return map[string]clicommand{
		"help": {
			name:        "help",
			description: "Displays a Help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Lists Next Location Areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists Previous Location Areas",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore {location_area}",
			description: "Lists the pokemon in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch {pokemon_area}",
			description: "Attempts to catch a pokemon and if caught adds it to pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect {pokemon_name}",
			description: "shows intersting info about pokemon, if caught",
			callback:    commandInspect,
		},

		"pokedex": {
			name:        "pokedex",
			description: "lists all the pokemons that are caught",
			callback:    commandPokedex,
		},
	}
}

func main() {

	cfg := config{
		pokeapiClient: pokeapi.NewClient(time.Hour),
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)

	command_exe := getCommands()
	for {
		fmt.Print("Pokedex>")
		scanner.Scan()
		text := scanner.Text()

		cleaned := cleanInput(text)

		if len(cleaned) == 0 {
			continue
		}
		commandName := cleaned[0]

		args := []string{}

		if len(cleaned) > 1 {
			args = cleaned[1:]
		}
		if _, ok := command_exe[commandName]; !ok {
			fmt.Println("Invalid Command")
			continue
		}
		err := command_exe[commandName].callback(&cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func cleanInput(str string) []string {
	lowered := strings.ToLower(str)
	words := strings.Fields(lowered)
	return words
}
