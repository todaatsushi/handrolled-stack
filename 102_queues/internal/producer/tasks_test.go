package producer_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/todaatsushi/queue/internal/messages"
	"github.com/todaatsushi/queue/internal/producer"
)

type writer struct {
	buffer *bytes.Buffer
}

func (w writer) Write(p []byte) (n int, err error) {
	w.buffer.Write(p)
	return 0, nil
}

type badWriter struct{}

func (w badWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("Some error")
}

func TestAddTasks(t *testing.T) {
	t.Run("Add task", func(t *testing.T) {
		var buf bytes.Buffer
		w := writer{buffer: &buf}
		msg := "Hello!"

		err := producer.QueueTask(w, msg)
		if err != nil {
			t.Fatal(err)
		}

		message := messages.NewMessage(messages.Enqueue, msg)
		expected, err := message.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		actual := w.buffer.Bytes()
		if !bytes.Equal(actual, expected) {
			t.Fatalf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Add task err handled", func(t *testing.T) {
		w := badWriter{}

		err := producer.QueueTask(w, "Hello!")
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := "Some error"
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}
