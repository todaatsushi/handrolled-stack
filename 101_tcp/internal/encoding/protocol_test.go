package encoding_test

import (
	"encoding/binary"
	"testing"

	"github.com/todaatsushi/basic_tcp/internal/encoding"
)

func TestEncoding(t *testing.T) {
	t.Run("Test headers", func(t *testing.T) {
		message := "Hello!"

		translator := encoding.Basic{}
		encoded := translator.Encode(message)

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
		message := "Hello!"

		translator := encoding.Basic{}
		encoded := translator.Encode(message)

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
