package client

import (
	"errors"
	"fmt"
	"strconv"
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

	key := parts[1]

	if command == protocol.Get {
		if len(parts) != 2 {
			return protocol.Message{}, errors.New("Invalid input, expected format: GET <key>.")
		}

		msg, err := protocol.NewMessage(command, key, []byte{}, -1)
		if err != nil {
			return protocol.Message{}, err
		}
		return msg, nil
	} else if command == protocol.Set {
		if len(parts) != 4 {
			return protocol.Message{}, errors.New("Invalid input, expected format: SET <key> <ttl> <data>.")
		}

		ttl, err := strconv.Atoi(parts[2])
		if err != nil {
			return protocol.Message{}, errors.New(fmt.Sprintf("Invalid input, couldn't parse '%s' to an int.", parts[2]))
		}

		data := parts[3]
		msg, err := protocol.NewMessage(command, key, []byte(data), ttl)
		if err != nil {
			return protocol.Message{}, err
		}
		return msg, nil
	} else {
		panic("unreachable")
	}
}
