package pokeapi

import (
	"bytes"
	"io"
	"os"
	"pokedex/internal/pokecache"
	"strings"
	"testing"
	"time"
)

func TestCatchAndAdd(t *testing.T) {
	cases := []struct {
		key string
		val bool
	}{
		{
			key: "pikachu",
			val: true,
		},
		{
			key: "eevee",
			val: true,
		},
		{
			key: "mew",
			val: true,
		},
	}

	for _, c := range cases {
		var p Pokedex
		p.PokemonCaught = make(map[string]Pokemon)
		var cfg Config
		cache := pokecache.NewCache(3 * time.Second)
		for i := 0; i < 50; i++ {
			Catch(&cfg, cache, &p, c.key)
			_, ok := p.PokemonCaught[c.key]
			if ok {
				break
			}
		}
		_, ok := p.PokemonCaught[c.key]
		if ok != c.val {
			t.Errorf("%s not found in pokedex", c.key)
			return
		}
	}
}

func TestShowPokedex(t *testing.T) {

	cases := []struct {
		key []Pokemon
		val bool
	}{
		{
			key: []Pokemon{
				{
					Name: "pikachu",
				},
				{
					Name: "Charizard",
				},
			},
			val: true,
		},
		{
			key: []Pokemon{
				{
					Name: "mewtwo",
				},
				{
					Name: "mew",
				},
				{
					Name: "entei",
				},
			},
			val: true,
		},
	}

	for _, c := range cases {
		cfg := Config{}
		cache := pokecache.NewCache(50 * time.Millisecond)
		pokedex := Pokedex{}
		pokedex.PokemonCaught = make(map[string]Pokemon)
		for _, poke := range c.key {
			pokedex.Add(poke)
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe")
		}
		os.Stdout = w

		ShowPokedex(&cfg, cache, &pokedex)

		w.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		if err != nil {
			t.Fatalf("failed to copy stdout.")
		}
		capturedOutput := buf.String()
		for _, poke := range c.key {
			contains := strings.Contains(capturedOutput, poke.Name)
			if contains != c.val {
				t.Errorf("Stdout does not contain '%s'", poke.Name)
			}
		}

		os.Stdout = oldStdout
	}

}
