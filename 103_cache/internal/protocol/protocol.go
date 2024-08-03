package protocol

import (
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

	panic("TODO")
}
