package protocol_test

import (
	"encoding/binary"
	"testing"

	"github.com/todaatsushi/handrolled-cache/internal/protocol"
)

type stream struct {
	data []byte
}

func (s stream) Read(buf []byte) (n int, err error) {
	copy(buf, s.data)
	return len(s.data), nil
}

func newStream(data []byte) stream {
	return stream{data}
}

func TestReader(t *testing.T) {
	t.Run("Reads SET", func(t *testing.T) {
		ttl := []byte{0, 1}

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte("hello,world!")
		dataLen := make([]byte, 2)
		binary.BigEndian.PutUint16(dataLen, uint16(len(data)))

		header := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		header = append(header, ttl...)
		header = append(header, keyLen...)
		header = append(header, dataLen...)

		message := []byte{}
		message = append(message, header...)
		message = append(message, key...)
		message = append(message, data...)

		s := newStream(message)
		reader := protocol.NewReader(s)

		read, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}

		if len(read) != len(message) {
			t.Errorf("Expected len %d, got %d", len(message), len(read))
		}

		for i, expected := range message {
			actual := read[i]
			if actual != expected {
				t.Errorf("Expected %d at position %d, got %d", expected, i, actual)
			}
		}
	})

	t.Run("Reads GET", func(t *testing.T) {
		ttl := []byte{0, 1}

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		dataLen := make([]byte, 2)

		header := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		header = append(header, ttl...)
		header = append(header, keyLen...)
		header = append(header, dataLen...)

		message := []byte{}
		message = append(message, header...)
		message = append(message, key...)

		s := newStream(message)
		reader := protocol.NewReader(s)

		read, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}

		if len(read) != len(message) {
			t.Errorf("Expected len %d, got %d", len(message), len(read))
		}

		for i, expected := range message {
			actual := read[i]
			if actual != expected {
				t.Errorf("Expected %d at position %d, got %d", expected, i, actual)
			}
		}
	})
}
