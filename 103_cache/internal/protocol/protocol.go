package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Version (1B) | Command (1B) | TTL (2B) | Length (2B) | Data (x)
const VERSION byte = 1
const HEADER_SIZE = 6

type Command byte

const (
	_ Command = iota
	Get
	Set
	Update
)

type Message struct {
	Cmd     Command
	Data    []byte
	Expires time.Time
}

func parseCommand(cmd byte) (Command, error) {
	switch cmd {
	case 1:
		return Get, nil
	case 2:
		return Set, nil
	case 3:
		return Update, nil
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
	}
	return nil
}

type Clock interface {
	Now() time.Time
	Add(d time.Duration) time.Time
}

func UnmarshalBinary(data []byte, clock Clock) (Message, error) {
	if len(data) < HEADER_SIZE {
		return Message{}, errors.New("Not enough data.")
	}

	version := data[0]
	if version != VERSION {
		return Message{}, errors.New("Version mismatch.")
	}

	lenDataBytes := data[4:6]
	lenData := int(binary.BigEndian.Uint16(lenDataBytes))
	var toCache []byte
	if lenData > 0 {
		toCache = data[6:]
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
		cmd, toCache, expires,
	}, nil
}
