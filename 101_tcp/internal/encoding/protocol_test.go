package encoding_test

import (
	"encoding/binary"
	"testing"

	"github.com/todaatsushi/basic_tcp/internal/encoding"
)

func TestEncoding(t *testing.T) {
	message := "Hello!"
	translator := encoding.Basic{}
	encoded := translator.Encode(message)

	t.Run("Test headers", func(t *testing.T) {
		actualVersion := encoded[0]
		if actualVersion != encoding.VERSION {
			t.Fatalf("Version flag mismatch: expected %d, got %d", encoding.VERSION, actualVersion)
		}

		expectedLenHeader := make([]byte, 2)
		binary.BigEndian.PutUint16(expectedLenHeader, uint16(len(message)))
		actualLenHeader := encoded[1:3]

		if len(actualLenHeader) != len(expectedLenHeader) {
			t.Fatalf("Expected len header to be of length %d, got %d", expectedLenHeader, actualLenHeader)
		}

		for i, expected := range expectedLenHeader {
			if expected != actualLenHeader[i] {
				t.Fatalf("Value mismatch of len header at index %d, expected %d, got %d", i, expected, actualLenHeader[i])
			}
		}
	})

	t.Run("Test data", func(t *testing.T) {
		expectedLen := 1 + 2 + len(message)
		if len(encoded) != expectedLen {
			t.Fatalf("Expected total size of %d, got %d", expectedLen, len(encoded))
		}

		actualData := encoded[encoding.HEADER_SIZE:]
		expectedData := []byte(message)

		if len(actualData) != len(expectedData) {
			t.Fatalf("Expected data to be %d long, got %d", len(expectedData), len(actualData))
		}

		for i, expected := range expectedData {
			if expected != actualData[i] {
				t.Fatalf("Value mismatch of len header at index %d, expected %d, got %d", i, expected, actualData[i])
			}
		}
	})
}

func TestDecoding(t *testing.T) {
	expected := "Hello!"
	translator := encoding.Basic{}

	t.Run("Decode", func(t *testing.T) {
		encoded := []byte{encoding.VERSION, 0, byte(len(expected))}
		encoded = append(encoded, []byte(expected)...)

		actual, err := translator.Decode(encoded)
		if err != nil {
			t.Fatalf("Error when decoding: %s", err.Error())
		}

		if actual != expected {
			t.Fatalf("Expected '%s', got '%s'.", expected, actual)
		}
	})

	t.Run("Not enough data", func(t *testing.T) {
		encoded := []byte{42}

		_, err := translator.Decode(encoded)
		if err == nil {
			t.Fatalf("Expected err, got nil.")
		}

		actual := err.Error()
		expectedErr := "Not enough data."
		if actual != expectedErr {
			t.Fatalf("Expected '%s', got '%s'", actual, expectedErr)
		}
	})

	t.Run("Version mismatch", func(t *testing.T) {
		encoded := []byte{42, 0, 0}

		_, err := translator.Decode(encoded)
		if err == nil {
			t.Fatalf("Expected err, got nil.")
		}

		expectedErr := "Version mismatch."
		actual := err.Error()
		if actual != expectedErr {
			t.Fatalf("Expected '%s', got '%s'", actual, expectedErr)
		}
	})

	t.Run("Data specified but no data", func(t *testing.T) {
		encoded := []byte{encoding.VERSION, 0, byte(len(expected))}

		_, err := translator.Decode(encoded)
		if err == nil {
			t.Fatalf("Expected err, got nil.")
		}

		expectedErr := "Data length specified but no data attached."
		actual := err.Error()
		if actual != expectedErr {
			t.Fatalf("Expected '%s', got '%s'", actual, expectedErr)
		}

	})

	t.Run("Len data mismatch with header", func(t *testing.T) {
		encoded := []byte{encoding.VERSION, 0, byte(len(expected))}
		encoded = append(encoded, byte(1))

		_, err := translator.Decode(encoded)
		if err == nil {
			t.Fatalf("Expected err, got nil.")
		}

		expectedErr := "Data size doesn't match length specified."
		actual := err.Error()
		if actual != expectedErr {
			t.Fatalf("Expected '%s', got '%s'", actual, expectedErr)
		}
	})
}
