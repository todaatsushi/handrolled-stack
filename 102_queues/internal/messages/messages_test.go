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
		expected = append(expected, '\n')

		actual, err := message.MarshalBinary()
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if !bytes.Equal(expected, actual) {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})

	t.Run("Marshal commands", func(t *testing.T) {
		testCases := []struct {
			Command  messages.Command
			Expected byte
		}{
			{
				messages.Log, 1,
			},
			{
				messages.Enqueue, 2,
			},
			{
				messages.Consume, 3,
			},
		}

		for _, tc := range testCases {
			message := messages.NewMessage(tc.Command, "")
			data, err := message.MarshalBinary()
			if err != nil {
				t.Error(err)
			}

			actual := data[1]
			if actual != tc.Expected {
				t.Errorf("Expected '%d', got '%d'", tc.Expected, actual)
			}
		}
	})

	t.Run("Consume message should have no data", func(t *testing.T) {
		message := messages.NewMessage(messages.Consume, "data")
		_, err := message.MarshalBinary()

		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := "Consume message should have no data."
		actual := err.Error()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("Unmarshal binary", func(t *testing.T) {
		data := []byte{
			1, // Version
			1, // Log
			0,
			1,          // Len of 1
			97,         // Data - 'a'
			byte('\n'), // Break
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
			1,          // Len of 1
			97,         // Data - 'a'
			byte('\n'), // Break
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
			1,          // Len of 1
			97,         // Data - 'a'
			byte('\n'), // Break
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
			byte('\n'), // Break
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
			byte('\n'), // Break
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
