package client_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/client"
	"github.com/todaatsushi/handrolled-cache/internal/protocol"
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
		cases := []struct {
			args     string
			expected string
		}{
			{"123", "Invalid input, expected format: SET <key> <ttl> <data>."},           // Just TTL
			{"", "Invalid input, expected format: SET <key> <ttl> <data>."},              // Not enough args
			{"InvalidTTL data", "Invalid input, couldn't parse 'InvalidTTL' to an int."}, // Invalid TTL
		}

		for _, tc := range cases {
			input := fmt.Sprintf("SET key %s", tc.args)

			_, err := client.ToMessage(input)
			if err == nil {
				t.Fatal("Expecting err, got nil.")
			}

			expected := errors.New(tc.expected).Error()
			actual := err.Error()

			if actual != expected {
				t.Errorf("Expected '%s', got '%s'", expected, actual)
			}
		}
	})

	t.Run("Test GET", func(t *testing.T) {
		input := "GET key"
		actual, err := client.ToMessage(input)
		if err != nil {
			t.Error(err)
		}

		expected := protocol.Message{
			Cmd:     protocol.Get,
			Key:     "key",
			Data:    []byte{},
			Expires: time.Now().Add(time.Second * time.Duration(420)),
		}

		if actual.Cmd != expected.Cmd {
			t.Errorf("Expected %d, got %d", expected.Cmd, actual.Cmd)
		}

		if actual.Key != expected.Key {
			t.Errorf("Expected %s, got %s", expected.Key, actual.Key)
		}

		// TODO - compare properly
		// if actual.Expires.Compare(expected.Expires) != 0 {
		// 	t.Errorf("Expected %s, got %s", expected.Expires, actual.Expires)
		// }

		if len(actual.Data) != len(expected.Data) {
			t.Errorf("Expecting %d, got %d", len(expected.Data), len(actual.Data))
		}

		for i, a := range actual.Data {
			e := expected.Data[i]
			if a != e {
				t.Errorf("Expected %d at position %d, got %d", e, i, a)
			}
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

			actual, err := client.ToMessage(input)
			if err != nil {
				t.Error(err)
			}

			expected := protocol.Message{
				Cmd:     protocol.Set,
				Key:     "key",
				Data:    []byte(tc),
				Expires: time.Now().Add(time.Second * time.Duration(420)),
			}

			if actual.Cmd != expected.Cmd {
				t.Errorf("Expected %d, got %d", expected.Cmd, actual.Cmd)
			}

			if actual.Key != expected.Key {
				t.Errorf("Expected %s, got %s", expected.Key, actual.Key)
			}

			// TODO - compare properly
			// if actual.Expires.Compare(expected.Expires) != 0 {
			// 	t.Errorf("Expected %s, got %s", expected.Expires, actual.Expires)
			// }

			if len(actual.Data) != len(expected.Data) {
				t.Errorf("Expecting %d, got %d", len(expected.Data), len(actual.Data))
			}

			for i, a := range actual.Data {
				e := expected.Data[i]
				if a != e {
					t.Errorf("Expected %d at position %d, got %d", e, i, a)
				}
			}
		}
	})
}
