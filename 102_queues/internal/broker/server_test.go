package broker_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/todaatsushi/queue/internal/broker"
	"github.com/todaatsushi/queue/internal/messages"
)

func TestHandle(t *testing.T) {
	t.Run("Log message", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)

		server := broker.NewServer(1337)
		message := messages.NewMessage(messages.Log, "Hello!")

		server.ProcessMessage(message)

		actual := buf.String()
		expected := "LOG: Hello!"

		if strings.Contains(actual, expected) == false {
			t.Errorf("Expected '%s' in log, got '%s'", expected, actual)
		}
	})

	t.Run("Enqueue message", func(t *testing.T) {
		server := broker.NewServer(1337)
		message := messages.NewMessage(messages.Enqueue, "Hello!")

		server.ProcessMessage(message)

		if server.QueueLen() != 1 {
			t.Errorf("Expected queue length of 1, got %v", server.QueueLen())
		}
	})

	t.Run("Consume message", func(t *testing.T) {
		t.Skip("TODO")
	})
}
