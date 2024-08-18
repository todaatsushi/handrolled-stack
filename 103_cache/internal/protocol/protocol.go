package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Version (1B) | Command (1B) | Expires (8B) | KeyLen (2B) | Length (2B) | Key (x) | Data (x)
const VERSION byte = 1
const HEADER_SIZE = 14

type Command byte

const (
	_ Command = iota
	Get
	Set
)

type Message struct {
	Cmd     Command
	Key     string
	Data    []byte
	Expires time.Time
}

func NewMessage(cmd Command, key string, data []byte, ttl int, c Clock) (Message, error) {
	if key == "" {
		return Message{}, errors.New("No key provided.")
	}

	if cmd == Set && len(data) == 0 {
		return Message{}, errors.New("No data provided for SET.")
	}

	if cmd == Get && len(data) > 0 {
		return Message{}, errors.New("Data provided for GET.")
	}

	if cmd == Get && ttl != 0 {
		return Message{}, errors.New("TTL must be 0 for GET.")
	}

	if cmd == Set && ttl <= 2 {
		return Message{}, errors.New("TTL must be greater than 2.")
	}

	expires := c.Now().Add(time.Second * time.Duration(ttl))
	return Message{cmd, key, data, expires}, nil
}

func (m Message) MarshalBinary(clock Clock) ([]byte, error) {
	expiresUnix := m.Expires.Unix()
	nowUnix := clock.Now().Unix()

	if m.Cmd == Set {
		if expiresUnix < nowUnix {
			return []byte{}, errors.New("Negative TTL.")
		}
	}

	expires := make([]byte, 8)
	binary.BigEndian.PutUint64(expires, uint64(expiresUnix))

	keyBytes := []byte(m.Key)
	keyLen := make([]byte, 2)
	binary.BigEndian.PutUint16(keyLen, uint16(len(keyBytes)))

	dataLen := make([]byte, 2)
	binary.BigEndian.PutUint16(dataLen, uint16(len(m.Data)))

	data := []byte{}
	data = append(data, VERSION)
	data = append(data, byte(m.Cmd))
	data = append(data, expires...)
	data = append(data, keyLen...)
	data = append(data, dataLen...)
	data = append(data, keyBytes...)
	data = append(data, m.Data...)
	return data, nil
}

func parseCommand(cmd byte) (Command, error) {
	switch cmd {
	case 1:
		return Get, nil
	case 2:
		return Set, nil
	default:
		return Get, errors.New(fmt.Sprintf("Invalid command: %d", int(cmd)))
	}
}

func validateData(cmd Command, data []byte, expires time.Time, clock Clock) error {
	switch cmd {
	case Get:
		if len(data) > 0 {
			return errors.New("Data passed to GET.")
		}
	case Set:
		if len(data) == 0 {
			return errors.New("Data not passed to SET.")
		}

		if expires.Compare(clock.Now()) < 0 {
			return errors.New("Expires in the past.")
		}
	}
	return nil
}

type Clock interface {
	Now() time.Time
}

func UnmarshalBinary(data []byte, clock Clock) (Message, error) {
	if len(data) < HEADER_SIZE {
		return Message{}, errors.New("Not enough data.")
	}

	version := data[0]
	if version != VERSION {
		return Message{}, errors.New("Version mismatch.")
	}

	keyLenBytes := data[10:12]
	lenKey := int(binary.BigEndian.Uint16(keyLenBytes))
	if lenKey == 0 {
		return Message{}, errors.New("No key provided.")
	}

	lenDataBytes := data[12:14]
	lenData := int(binary.BigEndian.Uint16(lenDataBytes))

	keyBytes := data[14 : 14+lenKey]
	key := string(keyBytes)

	var toCache []byte
	if lenData > 0 {
		toCache = data[14+lenKey:]
		if lenData != len(toCache) {
			return Message{}, errors.New("Length of data doesn't match header.")
		}
	} else {
		toCache = []byte{}
	}

	cmd, err := parseCommand(data[1])
	if err != nil {
		return Message{}, err
	}

	expiresBytes := data[2:10]
	expiresUnix := int64(binary.BigEndian.Uint64(expiresBytes))
	expires := time.Unix(expiresUnix, 0).UTC()

	err = validateData(cmd, toCache, expires, clock)
	if err != nil {
		return Message{}, err
	}

	return Message{
		cmd, key, toCache, expires,
	}, nil
}
