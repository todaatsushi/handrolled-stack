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

	log.Printf("Encoding '%s' and sending to server.", msg)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write(translator.Encode(msg)); err != nil {
		log.Fatal("Couldn't write encoded message: ", err)
	}
}
