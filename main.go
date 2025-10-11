package main

import (
	"strings"
)

func main() {
	cleanInput("   THIS isA      testCASE   ")
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	finalText := strings.Fields(lower)
	return finalText
}
