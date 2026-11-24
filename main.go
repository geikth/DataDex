package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const baseURL = "https://pokeapi.co/api/v2/pokemon"

type Liste struct {
	Results []PokemonItem `json:"results"`
}

type PokemonItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonDetails struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Height  int     `json:"height"`
	Weight  int     `json:"weight"`
	Sprites Sprites `json:"sprites"`
	Types   []Type  `json:"types"`
	Stats   []Stat  `json:"stats"`
}

type Sprites struct {
	FrontDefault string `json:"front_default"`
}

type Type struct {
	Type Resource `json:"type"`
}

type Stat struct {
	BaseStat int      `json:"base_stat"`
	Stat     Resource `json:"stat"`
}

type Resource struct {
	Name string `json:"name"`
}

func JSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("requête impossible : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("pokémon introuvable (erreur 404)")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu avec le code %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func PokemonList() ([]PokemonItem, error) {
	var list Liste
	if err := JSON(baseURL+"?limit=151&offset=0", &list); err != nil {
		return nil, err
	}
	return list.Results, nil
}

func Pokemon(name string) (*PokemonDetails, error) {
	var p PokemonDetails
	if err := JSON(baseURL+"/"+name, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func printList(pokemons []PokemonItem) {
	fmt.Println("\n=== POKÉDEX ===")
	for i, p := range pokemons {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}
}

func printDetails(p *PokemonDetails) {
	var types []string
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}
	stats := map[string]int{}
	for _, s := range p.Stats {
		stats[s.Stat.Name] = s.BaseStat
	}
	fmt.Printf("Nom        : %s\n", p.Name)
	fmt.Printf("Identifiant: %d\n", p.ID)
	fmt.Printf("Image      : %s\n", p.Sprites.FrontDefault)
	fmt.Printf("Taille     : %.1f m\n", float64(p.Height)/10)
	fmt.Printf("Poids      : %.1f kg\n", float64(p.Weight)/10)
	fmt.Printf("Types      : %s\n", strings.Join(types, ", "))
	fmt.Println("Stats :")
	fmt.Printf("  hp      : %d\n", stats["hp"])
	fmt.Printf("  attack  : %d\n", stats["attack"])
	fmt.Printf("  defense : %d\n", stats["defense"])
}

func main() {
	pokemons, err := PokemonList()
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		printList(pokemons)

		fmt.Print("\nNuméro du pokémon que tu veux:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > len(pokemons) {
			fmt.Printf("Entrée invalide : tape un nombre entre 1 et %d.\n", len(pokemons))
			continue
		}
		p, err := Pokemon(pokemons[n-1].Name)
		if err != nil {
			fmt.Println("Erreur :", err)
			continue
		}
		printDetails(p)

		fmt.Print("\nAppuie sur Entrée")
		reader.ReadString('\n')
	}
}
