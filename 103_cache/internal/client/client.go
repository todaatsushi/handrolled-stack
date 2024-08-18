package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/protocol"
)

type c struct{}

func (clock c) Now() time.Time {
	return time.Now().UTC()
}

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

		msg, err := protocol.NewMessage(command, key, []byte{}, 0, c{})
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
		msg, err := protocol.NewMessage(command, key, []byte(data), ttl, c{})
		if err != nil {
			return protocol.Message{}, err
		}
		return msg, nil
	} else {
		panic("unreachable")
	}
}

func Dial(port int) error {
	log.Println("Connecting client to port", port)
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		// TODO: 1 connection, multiple messages
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
		if err != nil {
			return err
		}

		go func() {
			// Read responses from server
			for scanner := bufio.NewScanner(conn); scanner.Scan(); {
				fmt.Printf("%s\n", scanner.Text())

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
		}()

		line := scanner.Text()

		msg, err := ToMessage(line)
		if err != nil {
			log.Println(err)
			continue
		}

		data, err := msg.MarshalBinary(c{})
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = conn.Write(data)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
