package pokepi

import (
	"fmt"
	"io"
	"net/http"
)

func Map() error {
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Println(body)
	return nil
}
