package producer

import (
	"log"
	"net"

	"github.com/todaatsushi/queue/internal/messages"
)

func CheckServer(port int) {
	message := messages.NewMessage(messages.Log, "healthcheck")
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
	} else {
		log.Println("Success.")
	}
}
