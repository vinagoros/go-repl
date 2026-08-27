package main

import "strings"

func cleanInput(text string) []string {
	trimmedText := strings.ToLower(strings.Trim(text, " "))
	return strings.Split(trimmedText, " ")
}
