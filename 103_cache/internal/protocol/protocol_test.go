package protocol_test

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/protocol"
)

type clock struct{}

func (c clock) Now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
	return t
}

func (c clock) Add(d time.Duration) time.Time {
	now := c.Now()
	return now.Add(d)
}

func TestUnmarshalValidation(t *testing.T) {
	t.Run("Not enough data", func(t *testing.T) {
		data := []byte{}

		expected := errors.New("Not enough data.").Error()
		_, err := protocol.UnmarshalBinary(data, clock{})

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

		_, err := protocol.UnmarshalBinary(data, clock{})
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

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Error("Expected err, got nil.")
		}

		expected := errors.New("Length of data doesn't match header.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Invalid command", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		data := []byte{
			protocol.VERSION,
			byte(0),
		}
		data = append(data, ttl...)
		data = append(data, size...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Error("Expected err, got nil.")
		}

		expected := errors.New("Invalid command: 0").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestUnmarshalGet(t *testing.T) {
	t.Run("Data passed to GET", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, size...)
		data = append(data, []byte{69}...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Error("Expected err, got nil.")
		}

		expected := errors.New("Data passed to GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}
