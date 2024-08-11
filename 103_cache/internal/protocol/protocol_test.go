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
			t.Fatal("Expected err, got nil.")
		}

		actual := err.Error()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'.", expected, actual)
		}
	})

	t.Run("No key", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		key := []byte("")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("No key provided.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Version mismatch", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)

		fakeVersion := byte(0)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			fakeVersion,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Version mismatch.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Data length mismatch", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 2)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Length of data doesn't match header.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Invalid command", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(0),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Invalid command: 0").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestUnmarshalGet(t *testing.T) {
	t.Run("Unmarshal GET", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)

		actual, err := protocol.UnmarshalBinary(data, clock{})
		if err != nil {
			t.Fatalf("Expected nil, got '%s'", err.Error())
		}

		expected := protocol.Message{
			protocol.Get, "key", []byte{}, clock{}.Now(),
		}

		if actual.Cmd != expected.Cmd {
			t.Errorf("Commands don't match: expected '%d', got '%d'", expected.Cmd, actual.Cmd)
		}
		if actual.Expires.Equal(expected.Expires) {
			t.Errorf("Expriry date doesn't match: expected '%s', got '%s'", expected.Expires, actual.Expires)
		}

		if len(actual.Data) != 0 {
			t.Errorf("No cached data expected for GET: Expected 0, got %d", len(actual.Data))
		}
	})

	t.Run("Data passed to GET", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Data passed to GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("TTL passed to GET", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("TTL shouldn't be passed to GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestUnmarshalSet(t *testing.T) {
	t.Run("Unmarshal SET", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		actual, err := protocol.UnmarshalBinary(data, clock{})
		if err != nil {
			t.Fatalf("Expected nil, got '%s'", err.Error())
		}

		expected := protocol.Message{
			protocol.Set, "key", []byte{69}, clock{}.Now(),
		}

		if actual.Cmd != expected.Cmd {
			t.Errorf("Commands don't match: expected '%d', got '%d'", expected.Cmd, actual.Cmd)
		}
		if actual.Expires.Equal(expected.Expires) {
			t.Errorf("Expriry date doesn't match: expected '%s', got '%s'", expected.Expires, actual.Expires)
		}

		if len(actual.Data) != 1 {
			t.Errorf("Cached data expected for SET: Expected 1, got %d", len(actual.Data))
		}
	})

	t.Run("Data not passed to SET", func(t *testing.T) {
		ttl := make([]byte, 2)
		secs := uint16(69)
		binary.BigEndian.PutUint16(ttl, secs)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Data not passed to SET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("TTL not passed to SET", func(t *testing.T) {
		ttl := make([]byte, 2)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, ttl...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("TTL not passed to SET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestMarshal(t *testing.T) {
	t.Run("Marshals", func(t *testing.T) {
		message := protocol.Message{
			protocol.Set, "key", []byte{69}, clock{}.Now().Add(time.Second * 10),
		}

		keyBytes := []byte("key")

		expected := []byte{}
		expected = append(expected, protocol.VERSION)
		expected = append(expected, byte(protocol.Set))
		expected = append(expected, []byte{0, 10}...)
		expected = append(expected, []byte{0, byte(len(keyBytes))}...)
		expected = append(expected, []byte{0, 1}...)
		expected = append(expected, keyBytes...)
		expected = append(expected, byte(69))

		actual, err := message.MarshalBinary(clock{})
		if err != nil {
			t.Fatal(err)
		}

		if len(actual) != len(expected) {
			t.Fatalf("Expected len %d, got %d", len(expected), len(actual))
		}

		for i, a := range actual {
			e := expected[i]

			if a != e {
				t.Errorf("Expected %d, got %d at position %d", e, a, i)
			}
		}
	})

	t.Run("Negative TTL", func(t *testing.T) {
		message := protocol.Message{
			protocol.Set, "key", []byte{69}, clock{}.Now().Add(time.Second * 10 * -1),
		}

		_, err := message.MarshalBinary(clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Negative TTL.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestNewMessage(t *testing.T) {
	t.Run("No key", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Get, "", []byte{}, 1)
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("No key provided.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("No data for SET", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Set, "key", []byte{}, 1)
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("No data provided for SET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Data provided for GET", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Get, "key", []byte{1}, 1)
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Data provided for GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("TTL less than 1 for SET", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Set, "key", []byte{1}, 0)
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("TTL must be greater than 0.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("TTL provided for GET", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Get, "key", []byte{}, 1)
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("TTL must be 0 for GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}
