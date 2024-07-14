package server

import (
	"bufio"
	"log"
	"net"

	"github.com/todaatsushi/basic_tcp/internal/encoding"
)

func Run(port *int, translator encoding.Translator) {
	log.Printf("Starting server on port %d.", *port)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(conn)
		scanner.Scan()

		decoded := translator.Decode(scanner.Bytes())
		log.Println("Received message: ", decoded)
	}
}
