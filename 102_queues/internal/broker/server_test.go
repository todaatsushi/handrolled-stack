package broker_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/todaatsushi/queue/internal/broker"
	"github.com/todaatsushi/queue/internal/messages"
)

type writer struct{}

func (w writer) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestHandle(t *testing.T) {
	t.Run("Log message", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)

		server := broker.NewServer(1337)
		message := messages.NewMessage(messages.Log, "Hello!")

		err := server.ProcessMessage(writer{}, message)
		if err != nil {
			t.Fatal(err)
		}

		actual := buf.String()
		expected := "LOG: Hello!"

		if strings.Contains(actual, expected) == false {
			t.Errorf("Expected '%s' in log, got '%s'", expected, actual)
		}
	})

	t.Run("Enqueue message", func(t *testing.T) {
		server := broker.NewServer(1337)
		message := messages.NewMessage(messages.Enqueue, "Hello!")

		err := server.ProcessMessage(writer{}, message)
		if err != nil {
			t.Fatal(err)
		}

		if server.QueueLen() != 1 {
			t.Errorf("Expected queue length of 1, got %v", server.QueueLen())
		}
	})

	t.Run("Consume message", func(t *testing.T) {
		t.Skip("TODO")
	})
}
