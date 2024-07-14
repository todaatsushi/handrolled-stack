package server

import (
	"log"
	"net"
)

func Run(port *int) {
	log.Printf("Starting server on port %d.", *port)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		_, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("TODO: implement binary parsing and log response.")
	}
}
