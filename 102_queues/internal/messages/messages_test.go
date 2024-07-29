package messages_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/todaatsushi/queue/internal/messages"
)

func TestMarshal(t *testing.T) {
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

func TestUnmarshal(t *testing.T) {
	t.Run("Unmarshal binary", func(t *testing.T) {
		data := []byte{
			1, // Version
			1, // Log
			0,
			1,  // Len of 1
			97, // Data - 'a'
		}

		expected := "a"

		actual, err := messages.UnmarshalBinary(data)
		if err != nil {
			t.Fatal(err)
		}

		if actual.Command != messages.Log {
			t.Errorf("Expected '%d', got '%d'", messages.Log, actual.Command)
		}

		if actual.Message != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual.Message)
		}
	})

	t.Run("Version mismatch", func(t *testing.T) {
		data := []byte{
			10, // Invalid version
			1,  // Log
			0,
			1,  // Len of 1
			97, // Data - 'a'
		}

		expected := errors.New("Version mismatch.")

		_, err := messages.UnmarshalBinary(data)
		if err == nil {
			t.Error("Expected error, got nil.")
		}

		if expected.Error() != err.Error() {
			t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
		}
	})

	t.Run("Invalid command", func(t *testing.T) {
		data := []byte{
			1,  // Version
			10, // Invalid command
			0,
			1,  // Len of 1
			97, // Data - 'a'
		}

		expected := errors.New("Unexpected command: 10")

		_, err := messages.UnmarshalBinary(data)
		if err == nil {
			t.Error("Expected error, got nil.")
		}

		if expected.Error() != err.Error() {
			t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
		}

	})

	t.Run("Message len doesn't match header", func(t *testing.T) {
		data := []byte{
			1, // Version
			1, // Command
			0,
			1,  // Len of 1
			97, // Len more than 1 data - 'aaa'
			97,
			97,
		}

		expected := errors.New("Mismatch in header info data length + received.")

		_, err := messages.UnmarshalBinary(data)
		if err == nil {
			t.Error("Expected error, got nil.")
		}

		if expected.Error() != err.Error() {
			t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
		}

	})

	t.Run("Message too short", func(t *testing.T) {
		data := []byte{
			1,
		}

		expected := errors.New("Not enough data.")

		_, err := messages.UnmarshalBinary(data)
		if err == nil {
			t.Error("Expected error, got nil.")
		}

		if expected.Error() != err.Error() {
			t.Errorf("Expected '%s', got '%s'", expected.Error(), err.Error())
		}

	})
}
