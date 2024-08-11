package client

import (
	"errors"
	"strings"

	"github.com/todaatsushi/handrolled-cache/internal/protocol"
)

func ToMessage(input string) (protocol.Message, error) {
	parts := strings.SplitN(input, " ", 4)
	if len(parts) < 2 {
		return protocol.Message{}, errors.New("Invalid format, should have 2/3 parts: CMD <KEY> <DATA (for SET)>")
	}

	cmd := strings.ToLower(parts[0])
	var command protocol.Command
	switch cmd {
	case "get":
		command = protocol.Get
	case "set":
		command = protocol.Set
	default:
		return protocol.Message{}, errors.New("Invalid command: should be SET or GET.")
	}

	if command == protocol.Get && len(parts) != 2 {
		return protocol.Message{}, errors.New("Invalid input, expected format: GET <key>.")
	} else if command == protocol.Set && len(parts) != 4 {
		return protocol.Message{}, errors.New("Invalid input, expected format: SET <key> <ttl> <data>.")
	}

	return protocol.Message{}, nil
}
