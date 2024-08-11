package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Version (1B) | Command (1B) | TTL (2B) | KeyLen (2B) | Length (2B) | Key (x) | Data (x)
const VERSION byte = 1
const HEADER_SIZE = 8

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

func NewMessage(cmd Command, key string, data []byte, ttl int) (Message, error) {
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

	if cmd == Set && ttl < 1 {
		return Message{}, errors.New("TTL must be greater than 0.")
	}

	expires := time.Now().Add(time.Second * time.Duration(ttl))
	return Message{cmd, key, data, expires}, nil
}

func (m Message) MarshalBinary(clock Clock) ([]byte, error) {
	diff := m.Expires.Sub(clock.Now())
	secs := diff.Seconds()

	if m.Cmd == Set {
		if secs < 0 {
			return []byte{}, errors.New("Negative TTL.")
		}
	}

	ttl := make([]byte, 2)
	binary.BigEndian.PutUint16(ttl, uint16(secs))

	keyBytes := []byte(m.Key)
	keyLen := make([]byte, 2)
	binary.BigEndian.PutUint16(keyLen, uint16(len(keyBytes)))

	dataLen := make([]byte, 2)
	binary.BigEndian.PutUint16(dataLen, uint16(len(m.Data)))

	data := []byte{}
	data = append(data, VERSION)
	data = append(data, byte(m.Cmd))
	data = append(data, ttl...)
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

func validateData(cmd Command, data []byte, ttl int) error {
	switch cmd {
	case Get:
		if len(data) > 0 {
			return errors.New("Data passed to GET.")
		}

		if ttl > 0 {
			return errors.New("TTL shouldn't be passed to GET.")
		}
	case Set:
		if len(data) == 0 {
			return errors.New("Data not passed to SET.")
		}

		if ttl == 0 {
			return errors.New("TTL not passed to SET.")
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

	keyLenBytes := data[4:6]
	lenKey := int(binary.BigEndian.Uint16(keyLenBytes))
	if lenKey == 0 {
		return Message{}, errors.New("No key provided.")
	}

	lenDataBytes := data[6:8]
	lenData := int(binary.BigEndian.Uint16(lenDataBytes))

	keyBytes := data[8 : 8+lenKey]
	key := string(keyBytes)

	var toCache []byte
	if lenData > 0 {
		toCache = data[8+lenKey:]
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

	TTLBytes := data[3:5]
	TTL := int(binary.BigEndian.Uint16(TTLBytes))
	expires := time.Now().Add(time.Second * time.Duration(TTL))

	err = validateData(cmd, toCache, TTL)
	if err != nil {
		return Message{}, err
	}

	return Message{
		cmd, key, toCache, expires,
	}, nil
}
