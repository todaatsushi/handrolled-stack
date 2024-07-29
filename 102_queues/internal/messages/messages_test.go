package messages_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/todaatsushi/queue/internal/messages"
)

func TestMessage(t *testing.T) {
	t.Run("Marshal binary from message", func(t *testing.T) {
		msg := "Hello, world!"
		message := messages.NewMessage(messages.Log, msg)

		expected := []byte{}
		expected = append(expected, messages.VERSION)
		expected = append(expected, byte(messages.Log))

		lenMessageData := make([]byte, 2)
		lenMessage := uint16(len(msg))
		binary.BigEndian.PutUint16(lenMessageData, lenMessage)

		expected = append(expected, lenMessageData...)
		expected = append(expected, []byte(msg)...)

		actual, err := message.MarshalBinary()
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if !bytes.Equal(expected, actual) {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})
}
