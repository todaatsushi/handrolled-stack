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
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Println(err.Error())
	}

	// Listen for response & write to log
	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		response := scanner.Bytes()
		// New scanner strips newline byte
		response = append(response, byte(messages.DELIM))

		message, err := messages.UnmarshalBinary(response)
		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Println("Num tasks:", message.Message)

		if err := scanner.Err(); err != nil {
			log.Println(err.Error())
		}
		break
	}
}
