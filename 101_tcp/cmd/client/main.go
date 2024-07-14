package client

import (
	"log"
	"net"
)

func encode(msg string) []byte {
	// Temp method - to be replaced
	log.Println("WARNING: not implemented")
	return []byte(msg)
}

func Send(msg string, port int) {
	if len(msg) == 0 {
		log.Fatal("Can't send empty message.")
	}

	log.Printf("Encoding '%s' and sending to server.", msg)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// TODO: impl encode message
	encoded := encode(msg)

	if _, err := conn.Write(encoded); err != nil {
		log.Fatal("Couldn't write encoded message: ", err)
	}
}
