package producer

import (
	"bufio"
	"log"
	"net"
	"time"

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
	defer conn.Close()

	// Listen for response & write to log
	go func() {
		for scanner := bufio.NewScanner(conn); scanner.Scan(); {

			response := scanner.Bytes()
			// New scanner strips newline byte
			response = append(response, byte(messages.DELIM))

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
	}
	// Stop closing connection before we can read out - make this less shit.
	time.Sleep(time.Second * 1)
}
