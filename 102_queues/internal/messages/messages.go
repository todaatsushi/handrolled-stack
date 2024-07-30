package messages

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Command int

const (
	_ Command = iota
	Log
	Enqueue
	Consume
)

const DELIM = '\n'

func parseCommand(value byte) (Command, error) {
	asInt := int(value)
	switch asInt {
	case 1:
		return Log, nil
	case 2:
		return Enqueue, nil
	case 3:
		return Consume, nil
	default:
		return -1, errors.New(fmt.Sprintf("Unexpected command: %d", asInt))
	}
}

// Header: Version (1B) | Command (1B) | LenMessage (2B)
const VERSION byte = 1

type Message struct {
	Command Command
	Message string
}

func NewMessage(command Command, message string) Message {
	return Message{
		Command: command,
		Message: message,
	}
}

func UnmarshalBinary(data []byte) (Message, error) {
	if len(data) <= 4 {
		return Message{}, errors.New("Not enough data.")
	}

	version := data[0]
	if version != VERSION {
		return Message{}, errors.New("Version mismatch.")
	}

	commandByte := data[1]
	lenMessageBytes := data[2:4]
	lenMessage := int(binary.BigEndian.Uint16(lenMessageBytes))

	// Header + data + break char
	if len(data) != 4+lenMessage+1 {
		return Message{}, errors.New("Mismatch in header info data length + received.")
	}

	command, err := parseCommand(commandByte)
	if err != nil {
		return Message{}, err
	}
	message := string(data[4 : 4+lenMessage])
	return NewMessage(command, message), nil
}

func (m Message) MarshalBinary() ([]byte, error) {
	data := []byte{}
	message := []byte(m.Message)

	var command byte
	switch m.Command {
	case Log:
		command = 1
	case Enqueue:
		command = 2
	case Consume:
		command = 3

		if len(m.Message) > 0 {
			return data, errors.New("Consume message should have no data.")
		}
	default:
		msg := fmt.Sprintf("Unhandled command: %d\n", m.Command)
		return data, errors.New(msg)
	}

	lenMessageData := make([]byte, 2)
	lenMessage := uint16(len(m.Message))
	binary.BigEndian.PutUint16(lenMessageData, lenMessage)

	data = append(data, VERSION)
	data = append(data, command)
	data = append(data, lenMessageData...)
	data = append(data, message...)
	data = append(data, DELIM)
	return data, nil
}
