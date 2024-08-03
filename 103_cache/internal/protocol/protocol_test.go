package protocol_test

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/todaatsushi/handrolled-cache/internal/protocol"
)

func TestUnmarshal(t *testing.T) {
	t.Run("Not enough data", func(t *testing.T) {
		data := []byte{}

		expected := errors.New("Not enough data.").Error()
		_, err := protocol.UnmarshalBinary(data)

		if err == nil {
			t.Error("Expected err, got nil.")
		}

		actual := err.Error()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'.", expected, actual)
		}
	})

	t.Run("Version mismatch", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		size := make([]byte, 2)

		fakeVersion := byte(0)

		data := []byte{
			fakeVersion,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, size...)

		_, err := protocol.UnmarshalBinary(data)
		if err == nil {
			t.Error("Expected err, got nil.")
		}

		expected := errors.New("Version mismatch.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Data length mismatch", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 2)

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, size...)
		data = append(data, []byte{69}...)

		_, err := protocol.UnmarshalBinary(data)
		if err == nil {
			t.Error("Expected err, got nil.")
		}

		expected := errors.New("Length of data doesn't match header.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}
