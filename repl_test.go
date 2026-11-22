package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	pokego "github.com/vinagoros/go-repl/internal"
)

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
		c := c
		t.Run(c.input, func(t *testing.T) {
			actual := cleanInput(c.input)

			if len(actual) != len(c.expected) {
				t.Fatalf("cleanInput(%q) returned %v; want %v", c.input, actual, c.expected)
			}
			for i := range actual {
				if actual[i] != c.expected[i] {
					t.Errorf("cleanInput(%q)[%d] = %q; want %q", c.input, i, actual[i], c.expected[i])
				}
			}
		})
	}
}

func TestCacheAddAndGet(t *testing.T) {
	cache := pokego.NewCache(time.Hour)
	cache.Add("testcase", []byte("test"))

	got, found := cache.Get("testcase")
	if !found {
		t.Fatal("Get() did not find the added entry")
	}
	if string(got) != "test" {
		t.Fatalf("Get() = %q; want %q", got, "test")
	}
	if _, found := cache.Get("missing"); found {
		t.Fatal("Get() found an entry that was never added")
	}
}

func TestCacheReapsExpiredEntry(t *testing.T) {
	const ttl = 40 * time.Millisecond
	cache := pokego.NewCache(ttl)
	cache.Add("expired", []byte("value"))

	time.Sleep(3 * ttl)
	if got, found := cache.Get("expired"); found {
		t.Fatalf("Get() found expired entry %q", got)
	}
}

func TestPokeGetUsesCache(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "test"})
	}))
	defer server.Close()

	for i := 0; i < 6; i++ {
		got, err := pokego.PokeGet[map[string]string](server.URL)
		if err != nil {
			t.Fatalf("PokeGet() error: %v", err)
		}
		if got["name"] != "test" {
			t.Fatalf("got name %q, want %q", got["name"], "test")
		}
	}

	if got := requests.Load(); got != 1 {
		t.Fatalf("server received %d requests, want 1", got)
	}
}

func TestPokeGetReturnsErrorForNonOKResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	if _, err := pokego.PokeGet[map[string]string](server.URL); err == nil {
		t.Fatal("PokeGet() returned no error for a non-200 response")
	}
}

func TestPokeGetReturnsErrorForInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("{"))
	}))
	defer server.Close()

	if _, err := pokego.PokeGet[map[string]string](server.URL); err == nil {
		t.Fatal("PokeGet() returned no error for invalid JSON")
	}
}
