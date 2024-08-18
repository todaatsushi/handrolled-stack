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
	t, _ := time.Parse(time.RFC3339, "2069-04-20T15:00:00Z")
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
		expires := make([]byte, 8)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		key := []byte("")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, expires...)
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
		expires := make([]byte, 8)

		size := make([]byte, 2)

		fakeVersion := byte(0)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			fakeVersion,
			byte(protocol.Get),
		}
		data = append(data, expires...)
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
		expires := make([]byte, 8)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 2)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, expires...)
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
		ttl := make([]byte, 8)

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
		expires := make([]byte, 8)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, expires...)
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
		expires := make([]byte, 8)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Get),
		}
		data = append(data, expires...)
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
}

func TestUnmarshalSet(t *testing.T) {
	t.Run("Unmarshal SET", func(t *testing.T) {
		expiresAt := clock{}.Now().Add(time.Second * time.Duration(69)).UTC()

		expires := make([]byte, 8)
		unix := uint64(expiresAt.Unix())
		binary.BigEndian.PutUint64(expires, unix)

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, expires...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		actual, err := protocol.UnmarshalBinary(data, clock{})
		if err != nil {
			t.Fatalf("Expected nil, got '%s'", err.Error())
		}

		expected := protocol.Message{
			protocol.Set, "key", []byte{69}, expiresAt,
		}

		if actual.Cmd != expected.Cmd {
			t.Errorf("Commands don't match: expected '%d', got '%d'", expected.Cmd, actual.Cmd)
		}
		if !actual.Expires.Equal(expected.Expires) {
			t.Errorf("Expriry date doesn't match: expected '%s', got '%s'", expected.Expires, actual.Expires)
		}

		if len(actual.Data) != 1 {
			t.Errorf("Cached data expected for SET: Expected 1, got %d", len(actual.Data))
		}
	})

	t.Run("Data not passed to SET", func(t *testing.T) {
		expires := make([]byte, 8)
		secs := uint64(clock{}.Now().Add(time.Second * time.Duration(69)).Unix())
		binary.BigEndian.PutUint64(expires, secs)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 0)

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, expires...)
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

	t.Run("Expires in the past", func(t *testing.T) {
		expires := make([]byte, 8)
		old := clock{}.Now().Add(time.Second * time.Duration(-3600))
		binary.BigEndian.PutUint64(expires, uint64(old.Unix()))

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, expires...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		_, err := protocol.UnmarshalBinary(data, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Expires in the past.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Expires parsed correctly", func(t *testing.T) {
		expires := make([]byte, 8)
		dt := clock{}.Now().Add(time.Second * time.Duration(1800))
		binary.BigEndian.PutUint64(expires, uint64(dt.Unix()))

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, 1)

		key := []byte("key")
		keyLen := make([]byte, 2)
		binary.BigEndian.PutUint16(keyLen, uint16(len(key)))

		data := []byte{
			protocol.VERSION,
			byte(protocol.Set),
		}
		data = append(data, expires...)
		data = append(data, keyLen...)
		data = append(data, size...)
		data = append(data, key...)
		data = append(data, byte(69))

		m, err := protocol.UnmarshalBinary(data, clock{})
		if err != nil {
			t.Fatal(err)
		}

		expected := dt.Unix()
		actual := m.Expires.Unix()

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})
}

func TestMarshal(t *testing.T) {
	t.Run("Marshals", func(t *testing.T) {
		expires := clock{}.Now().Add(time.Second * 10)
		message := protocol.Message{
			protocol.Set, "key", []byte{69}, clock{}.Now().Add(time.Second * 10),
		}

		keyBytes := []byte("key")

		expiresBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(expiresBytes, uint64(expires.Unix()))

		expected := []byte{}
		expected = append(expected, protocol.VERSION)
		expected = append(expected, byte(protocol.Set))
		expected = append(expected, expiresBytes...)
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

	t.Run("Expires marshalled correctly", func(t *testing.T) {
		expectedDt := clock{}.Now().Add(time.Second * 10)
		expectedUnix := int(expectedDt.Unix())

		expectedBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedBytes, uint64(expectedUnix))

		message := protocol.Message{
			protocol.Set, "key", []byte{69}, expectedDt,
		}

		data, err := message.MarshalBinary(clock{})
		if err != nil {
			t.Fatal(err)
		}

		expiresBytes := data[2:10]
		actualUnix := int(binary.BigEndian.Uint64(expiresBytes))

		if actualUnix != expectedUnix {
			t.Errorf("Expected %d, got %d", expectedUnix, actualUnix)
		}

		actual := time.Unix(int64(actualUnix), 0).UTC()
		if expectedDt.Compare(actual) != 0 {
			t.Errorf("Expected '%s', got '%s'", expectedDt, actual)
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
		_, err := protocol.NewMessage(protocol.Get, "", []byte{}, 1, clock{})
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
		_, err := protocol.NewMessage(protocol.Set, "key", []byte{}, 10, clock{})
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
		_, err := protocol.NewMessage(protocol.Get, "key", []byte{1}, 1, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("Data provided for GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("TTL less than 3 for SET", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Set, "key", []byte{1}, 0, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("TTL must be greater than 2.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("TTL provided for GET", func(t *testing.T) {
		_, err := protocol.NewMessage(protocol.Get, "key", []byte{}, 1, clock{})
		if err == nil {
			t.Fatal("Expected err, got nil.")
		}

		expected := errors.New("TTL must be 0 for GET.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("New message", func(t *testing.T) {
		msg, err := protocol.NewMessage(protocol.Set, "key", []byte{1}, 10, clock{})
		if err != nil {
			t.Fatal(err)
		}

		actual := msg.Expires
		fakeCurrent := clock{}.Now()

		expectedYear := fakeCurrent.Year()
		expectedDay := fakeCurrent.Day()
		expectedMonth := fakeCurrent.Month()
		expectedHour := fakeCurrent.Hour()
		expectedMinute := fakeCurrent.Minute()
		expectedSecond := fakeCurrent.Second() + 10

		if actual.Year() != expectedYear {
			t.Errorf("Expected year %d, got %d", expectedYear, actual.Year())
		}

		if actual.Month() != expectedMonth {
			t.Errorf("Expected month %d, got %d", expectedMonth, actual.Month())
		}

		if actual.Day() != expectedDay {
			t.Errorf("Expected day %d, got %d", expectedDay, actual.Day())
		}

		if actual.Hour() != expectedHour {
			t.Errorf("Expected hour %d, got %d", expectedHour, actual.Hour())
		}

		if actual.Minute() != expectedMinute {
			t.Errorf("Expected minute %d, got %d", expectedMinute, actual.Minute())
		}

		if actual.Second() != expectedSecond {
			t.Errorf("Expected second %d, got %d", expectedSecond, actual.Second())
		}
	})
}

func TestEncodeDecode(t *testing.T) {
	t.Run("Marshall / unmarshall SET", func(t *testing.T) {
		key := "key"
		data := []byte("Some data")

		c := clock{}
		expires, _ := time.Parse(time.RFC3339, "2069-04-20T15:30:00Z")
		ttl := int(expires.Sub(c.Now()).Seconds())

		if ttl != 1800 {
			t.Fatalf("TTL should be 1800, not %d", ttl)
		}

		expected, err := protocol.NewMessage(
			protocol.Set, key, data, ttl, c,
		)
		if err != nil {
			t.Fatal(err)
		}

		encoded, err := expected.MarshalBinary(c)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := protocol.UnmarshalBinary(encoded, c)
		if err != nil {
			t.Fatal(err)
		}

		if actual.Cmd != expected.Cmd {
			t.Errorf("Expected cmd %d, got %d", expected.Cmd, actual.Cmd)
		}

		if actual.Key != expected.Key {
			t.Errorf("Expected key '%s', got '%s'", expected.Key, actual.Key)
		}

		if actual.Expires.Compare(expected.Expires) != 0 {
			t.Errorf("Expected expires '%s', got '%s'", expected.Expires, actual.Expires)
		}

		if actual.Expires.Compare(expires) != 0 {
			t.Errorf("Expected expires '%s', got '%s'", expires, actual.Expires)
		}

		if len(actual.Data) != len(expected.Data) {
			t.Errorf("Expected len data %d, got %d", len(expected.Data), len(actual.Data))
		}
	})
}
