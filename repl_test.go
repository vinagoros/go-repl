package main

import "testing"

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    " three word sentence ",
			expected: []string{"three", "word", "sentence"},
		},
		{
			input:    "o rato ROEU a rOlHa da garrafa",
			expected: []string{"o", "rato", "roeu", "a", "rolha", "da", "garrafa"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected: %v, Actual: %v", c.expected, actual)
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected: %s, Actual: %s", expectedWord, word)
				t.Fail()
			}
		}
	}
}
