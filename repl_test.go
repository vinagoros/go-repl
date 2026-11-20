package main

import (
	"encoding/json"
	"fmt"
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

type input struct {
	key   string
	value []byte
}
type expected struct {
	value []byte
	found bool
}

// ensure cache is properly adding,
func TestAddToCache(t *testing.T) {
	cases := []struct {
		input    input
		expected expected
	}{
		{
			input{
				key:   "testcase",
				value: []byte("test"),
			},
			expected{
				value: []byte("test"),
				found: true,
			},
		},
	}
	fmt.Println("Cache Get Test")
	testCache := pokego.NewCache(60)
	fmt.Println("Running test cases")
	for _, testCase := range cases {
		entryKey := testCase.input.key
		entryValue := testCase.input.value
		testCache.Add(entryKey, entryValue)
		result, found := testCache.Get(entryKey)
		fmt.Printf("Input: %s %s\nExpected: %s %t\nResult: %s %t\n", entryKey, entryValue,
			entryValue, true,
			result, found)
	}

}

func TestReapLoop(T *testing.T) {
	cases := []struct {
		input    input
		expected expected
	}{
		{
			input{
				key:   "timetestcase",
				value: []byte("timehasnotpassedyet"),
			},
			expected{
				value: nil,
				found: false,
			},
		},
	}
	fmt.Println("Cache Reaploop Test")
	testCache := pokego.NewCache(5 * time.Second)
	fmt.Println("Running test cases")
	for _, testCase := range cases {
		entryKey := testCase.input.key
		entryValue := testCase.input.value
		testCache.Add(entryKey, entryValue)
		time.Sleep(6 * time.Second)
		result, found := testCache.Get(entryKey)
		fmt.Printf("Input: %s %s\nExpected: %s %t\nResult: %s %t\n", entryKey, entryValue,
			testCase.expected.value, false,
			result, found)
	}

}

func TestPokeGetUsesCache(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "test"})
	}))
	defer server.Close()

	for i := 0; i < 2; i++ {
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
