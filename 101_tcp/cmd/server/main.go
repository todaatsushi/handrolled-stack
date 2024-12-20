package server

import (
	"io"
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
		log.Println("Connected.")

		go handle(conn, translator)
	}
}

func handle(conn net.Conn, translator encoding.Translator) {
	r := encoding.NewMessageReader(conn)

	for {
		data, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
		msg, err := translator.Decode(data)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Received message: ", msg)
	}
}
