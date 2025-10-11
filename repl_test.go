package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "This isatest here",
			expected: []string{"this", "isatest", "here"},
		},
		{
			input:    "Esther Forseth         rules     ",
			expected: []string{"esther", "forseth", "rules"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Test failed: Resulting slice length doesn't match expected length")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Test failed: Resulting words do not match expected words")
			}
		}
	}
}
