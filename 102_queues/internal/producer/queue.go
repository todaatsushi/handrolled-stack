package producer

import (
	"bufio"
	"log"
	"net"

	"github.com/todaatsushi/queue/internal/messages"
)

func GetQueueLen(port int) {
	message := messages.NewMessage(messages.QueueLen, "")
	data, err := message.MarshalBinary()
	if err != nil {
		log.Println(err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Listen for response & write to log
	go func() {
		for scanner := bufio.NewScanner(conn); scanner.Scan(); {
			response := scanner.Bytes()
			message, err := messages.UnmarshalBinary(response)
			if err != nil {
				log.Println(err.Error())
				return
			}

			log.Println(message.Message)

			if err := scanner.Err(); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	_, err = conn.Write(data)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Success.")
	}
}
