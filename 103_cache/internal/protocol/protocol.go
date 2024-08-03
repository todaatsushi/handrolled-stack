package protocol

import (
	"encoding/binary"
	"errors"
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

func UnmarshalBinary(data []byte) (Message, error) {
	if len(data) < HEADER_SIZE {
		return Message{}, errors.New("Not enough data.")
	}

	version := data[0]
	if version != VERSION {
		return Message{}, errors.New("Version mismatch.")
	}

	lenDataBytes := data[4:6]
	lenData := int(binary.BigEndian.Uint16(lenDataBytes))
	toCache := data[7:]
	if lenData != len(toCache) {
		return Message{}, errors.New("Length of data doesn't match header.")
	}

	panic("TODO")
}
