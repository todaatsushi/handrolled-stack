package consumer

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/todaatsushi/queue/internal/messages"
)

func processMessage(conn net.Conn) error {
	message := messages.NewMessage(messages.Consume, "")
	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		response := scanner.Bytes()
		// New scanner strips newline byte
		response = append(response, byte(messages.DELIM))

		msg, err := messages.UnmarshalBinary(response)
		if err != nil {
			return err
		}
		if msg.Command == messages.Consume {
			// TODO: replace with new command to signify no task
			return nil
		}
		log.Println("Message:", msg)
		break
	}
	return nil
}

func poll(port int) error {
	for {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
		if err != nil {
			return err
		}

		err = processMessage(conn)
		if err != nil {
			return err
		}

		// Fake processing the message.
		interval := rand.Intn(3)
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func StartConsumers(port int, numConsumers int) {
	errs := make(chan error)
	for range numConsumers {
		go func(p int) {
			err := poll(p)
			if err != nil {
				errs <- err
			}
		}(port)
	}
	<-errs
}
