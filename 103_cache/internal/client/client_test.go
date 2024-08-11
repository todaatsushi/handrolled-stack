package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/todaatsushi/handrolled-cache/internal/client"
)

func TestSend(t *testing.T) {
	t.Run("Not enough args", func(t *testing.T) {
		input := ""

		_, err := client.ToMessage(input)
		if err == nil {
			t.Fatal("Expecting err, got nil.")
		}

		expected := errors.New("Invalid format, should have 2/3 parts: CMD <KEY> <DATA (for SET)>").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Invalid command", func(t *testing.T) {
		input := "TEST hello"

		_, err := client.ToMessage(input)
		if err == nil {
			t.Fatal("Expecting err, got nil.")
		}

		expected := errors.New("Invalid command: should be SET or GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Invalid GET", func(t *testing.T) {
		input := "GET hello 123"

		_, err := client.ToMessage(input)
		if err == nil {
			t.Fatal("Expecting err, got nil.")
		}

		expected := errors.New("Invalid input, expected format: GET <key>.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Invalid SET", func(t *testing.T) {
		cases := []string{
			"123",        // Just TTL
			"InvalidTTL", // TTL should be int
			"",           // No extras
		}

		for _, tc := range cases {
			input := fmt.Sprintf("SET key %s", tc)

			_, err := client.ToMessage(input)
			if err == nil {
				t.Fatal("Expecting err, got nil.")
			}

			expected := errors.New("Invalid input, expected format: SET <key> <ttl> <data>.").Error()
			actual := err.Error()

			if actual != expected {
				t.Errorf("Expected '%s', got '%s'", expected, actual)
			}
		}
	})

	t.Run("Test GET", func(t *testing.T) {
		input := "GET key"
		_, err := client.ToMessage(input)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test SET", func(t *testing.T) {
		cases := []string{
			"data",
			"multiple data input with spaces",
			"{'fake': 'json', 'nested': {'json': true}}",
		}

		for _, tc := range cases {
			input := fmt.Sprintf("SET key 420 %s", tc)

			_, err := client.ToMessage(input)
			if err != nil {
				t.Error(err)
			}
		}
	})
}
