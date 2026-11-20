package pokego

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var cache = NewCache(
	10 * time.Minute,
)

func PokeGet[T any](endpoint string) (T, error) {
	var data T
	// TO DO passar a reconhecer o tipo de struct
	var output []byte
	result, found := cache.Get(endpoint)
	var err error
	if found {
		output = result
	} else {
		resp, err := http.Get(endpoint)
		if err != nil {
			return data, fmt.Errorf("Unable to retrieve data: %w", err)
		}
		if resp.StatusCode != 200 {
			return data, fmt.Errorf("Error in query: %s", resp.Status)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return data, fmt.Errorf("Error while reading response: %w", err)
		}
		output = body
	}
	err = json.Unmarshal(output, &data)
	if err != nil {
		return data, fmt.Errorf("Error while parsing to json response: %w", err)
	}
	cache.Add(endpoint, output)

	return data, nil

}
