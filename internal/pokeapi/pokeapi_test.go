package pokeapi

import (
	"pokedex/internal/pokecache"
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
