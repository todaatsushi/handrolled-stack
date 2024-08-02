package protocol_test

import (
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
}
