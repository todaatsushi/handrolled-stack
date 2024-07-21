package client

import (
	"log"
	"net"

	"github.com/todaatsushi/basic_tcp/internal/encoding"
)

func Send(msg string, port int, translator encoding.Translator) {
	if len(msg) == 0 {
		log.Fatal("Can't send empty message.")
	}

	log.Printf("Encoding '%s' and sending to server (%d).", msg, len(msg))

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Artificial stream creation - it should send as one message
	encoded := translator.Encode(msg)

	for _, b := range encoded {
		if _, err := conn.Write([]byte{b}); err != nil {
			log.Fatal(err)
		}
	}
}
