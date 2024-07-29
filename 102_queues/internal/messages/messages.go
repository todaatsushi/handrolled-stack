package messages

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Command int

const (
	_           = iota
	Log Command = 1
)

// Header: Version (1B) | Command (1B) | LenMessage (2B)
const VERSION byte = 1

type Message struct {
	command Command
	message string
}

func NewMessage(command Command, message string) Message {
	return Message{
		command: command,
		message: message,
	}
}

func (m Message) UnmarshalBinary(data []byte) Message {
	return Message{
		command: 1,
		message: "TODO",
	}
}

func (m Message) MarshalBinary() ([]byte, error) {
	data := []byte{}
	message := []byte(m.message)

	var command byte
	switch m.command {
	case Log:
		command = 1
	default:
		msg := fmt.Sprintf("Unhandled command: %d\n", m.command)
		return data, errors.New(msg)
	}

	lenMessageData := make([]byte, 2)
	lenMessage := uint16(len(m.message))
	binary.BigEndian.PutUint16(lenMessageData, lenMessage)

	data = append(data, VERSION)
	data = append(data, command)
	data = append(data, lenMessageData...)
	data = append(data, message...)
	return data, nil
}
