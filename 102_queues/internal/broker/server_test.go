package broker_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/todaatsushi/queue/internal/broker"
	"github.com/todaatsushi/queue/internal/messages"
)

type writer struct {
	buffer *bytes.Buffer
}

func (w writer) Write(p []byte) (n int, err error) {
	w.buffer.Write(p)
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
		server := broker.NewServer(1337)
		data := "Hello!"
		expected := messages.NewMessage(messages.Enqueue, data)

		var buf bytes.Buffer
		w := writer{
			buffer: &buf,
		}

		err := server.ProcessMessage(w, expected)
		if err != nil {
			t.Fatal(err)
		}

		if server.QueueLen() != 1 {
			t.Errorf("Expected queue length to be 1, got %d", server.QueueLen())
		}

		message := messages.NewMessage(messages.Consume, "")
		err = server.ProcessMessage(w, message)
		if err != nil {
			t.Fatal(err)
		}

		writtenMessage := w.buffer.Bytes()
		parsedMessage, err := messages.UnmarshalBinary(writtenMessage)
		if err != nil {
			t.Fatal(err)
		}
		actual := parsedMessage.Message
		if actual != data {
			t.Errorf("Expected '%s', got '%s'", data, actual)
		}

		if server.QueueLen() != 0 {
			t.Errorf("Expected queue length to be 0, got %d", server.QueueLen())
		}
	})

	t.Run("Get queue length", func(t *testing.T) {
		server := broker.NewServer(1337)
		expected := messages.NewMessage(messages.Enqueue, "Hello!")

		var buf bytes.Buffer
		w := writer{
			buffer: &buf,
		}

		err := server.ProcessMessage(w, expected)
		if err != nil {
			t.Fatal(err)
		}

		if server.QueueLen() != 1 {
			t.Errorf("Expected queue length to be 1, got %d", server.QueueLen())
		}

		message := messages.NewMessage(messages.QueueLen, "")
		err = server.ProcessMessage(w, message)
		if err != nil {
			t.Fatal(err)
		}

		writtenMessage := w.buffer.Bytes()
		parsedMessage, err := messages.UnmarshalBinary(writtenMessage)
		if err != nil {
			t.Fatal(err)
		}
		actual := parsedMessage.Message
		if actual != "1" {
			t.Errorf("Expected '1', got '%s'", actual)
		}
	})
}
